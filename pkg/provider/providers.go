package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/axatol/home-assistant-integrations/pkg/broker"
	"github.com/axatol/home-assistant-integrations/pkg/homeassistant"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Provider interface {
	// Name returns the name of the provider.
	Name() string
	// Configure gives the provider a chance to initialise and configure clients.
	Configure() error
	// StateTopic is the topic that device state will be published to.
	StateTopic() string
	// AvailabilityTopic is the topic that individual device availability status will be published.
	AvailabilityTopic() string
	// Interval provides a ticker to schedule poll attempts.
	Interval() <-chan time.Time
	// DeviceMetadata returns metadata about the device.
	DeviceMetadata() homeassistant.DeviceInformation
	// EntityConfigurationSet returns the configuration for the entities provided by the device to announce to Home Assistant.
	EntityConfigurationSet() map[string]homeassistant.EntityConfiguration
	// Poll interrogates the device state
	Poll(ctx context.Context) (map[string]any, error)
	// Health checks for errors and responds with provider metadata.
	Health(ctx context.Context) (map[string]any, error)
}

var Providers = []Provider{
	new(HuaweiHG659Provider),
	new(ZeverSolarTLC5000Provider),
}

func Configure(ctx context.Context) error {
	for _, provider := range Providers {
		name := provider.Name()
		log.Debug().Msgf("configuring provider %s", name)

		if err := provider.Configure(); err != nil {
			return fmt.Errorf("failed to configure provider %s: %s", name, err)
		}

		schema := provider.EntityConfigurationSet()
		if schema == nil {
			log.Info().Msgf("skipping disabled provider %s", name)
			continue
		}

		Providers = append(Providers, provider)
	}

	return nil
}

func Start(ctx context.Context, mqtt *broker.Broker) error {
	g := errgroup.Group{}

	for _, provider := range Providers {
		p := provider
		g.Go(func() error {
			for topic, entity := range p.EntityConfigurationSet() {
				if err := mqtt.Publish(ctx, topic, entity, broker.WithRetained(true)); err != nil {
					return fmt.Errorf("failed to announce %s/%s: %s", p.Name(), entity.Name, err)
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	for _, provider := range Providers {
		p := provider
		go func() {
			log := log.With().Str("provider", p.Name()).Logger()
			ctx = log.WithContext(ctx)

			for {
				select {
				case <-ctx.Done():
					return

				case <-p.Interval():
					data, err := p.Poll(ctx)
					if err != nil {
						log.Error().Err(fmt.Errorf("failed to poll: %s", err)).Send()
						mqtt.Publish(ctx, p.AvailabilityTopic(), "offline")
						continue
					}

					if data == nil {
						mqtt.Publish(ctx, p.AvailabilityTopic(), "offline")
						continue
					}

					if err := mqtt.Publish(ctx, p.StateTopic(), data); err != nil {
						log.Error().Err(fmt.Errorf("failed to publish data: %s", err)).Send()
					}
				}
			}
		}()
	}

	return g.Wait()
}
