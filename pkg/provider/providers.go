package provider

import (
	"context"
	"fmt"

	"github.com/axatol/home-assistant-integrations/pkg/broker"
	"github.com/axatol/home-assistant-integrations/pkg/config"
	"github.com/axatol/home-assistant-integrations/pkg/homeassistant"
	"github.com/rs/zerolog/log"
)

const (
	PROVIDER_SINK_BUFFER_SIZE = 50
)

type Provider interface {
	// Name returns the name of the provider.
	Name() string
	// Configure gives the provider a chance to initialise and configure clients.
	Configure() error
	// Schema returns the entity configuration to announce to Home Assistant.
	Schema() map[string]homeassistant.EntityConfiguration
	// Health checks for errors and responds with provider metadata.
	Health(ctx context.Context) (map[string]any, error)
	// Subscribe produces a channel for the manager to listen to.
	Subscribe(ctx context.Context) <-chan broker.Payload
}

var (
	deviceOrigin = homeassistant.EntityOrigin{
		Name:            "Home Assistant Integrations",
		SupportURL:      "https://github.com/axatol/home-assistant-integrations/issues",
		SoftwareVersion: config.BuildVersion,
	}

	AvailableProviders = []Provider{
		new(HuaweiHG659Provider),
		new(ZeverSolarTLC5000Provider),
	}

	Providers = []Provider{}
)

func Configure(ctx context.Context) error {
	for _, provider := range AvailableProviders {
		name := provider.Name()
		log.Debug().Msgf("configuring provider %s", name)

		if err := provider.Configure(); err != nil {
			return fmt.Errorf("failed to configure provider %s: %s", name, err)
		}

		schema := provider.Schema()
		if schema == nil {
			log.Info().Msgf("skipping disabled provider %s", name)
			continue
		}

		Providers = append(Providers, provider)
	}

	return nil
}
