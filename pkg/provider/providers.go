package provider

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Provider interface {
	// Name returns the name of the provider.
	Name() string
	// Configure initialises and configures with options.
	Configure(options map[string]any) error
	// Health checks for errors and responds with provider metadata.
	Health(ctx context.Context) (map[string]any, error)
	// Close closes the provider.
	Close() error
}

var (
	Providers = map[string]Provider{
		"advantage_air_hub": &AdvantageAirHubProvider{},
	}

	ConfiguredProviders = map[string]Provider{}
)

func Configure(ctx context.Context, configs map[string]map[string]any) error {
	for name, config := range configs {
		log.Debug().Msgf("configuring provider %s", name)

		provider, ok := Providers[name]
		if !ok {
			log.Warn().Err(fmt.Errorf("unknown provider: %s", name)).Send()
			continue
		}

		if err := provider.Configure(config); err != nil {
			return fmt.Errorf("failed to configure provider %s: %s", name, err)
		}

		if _, err := provider.Health(ctx); err != nil {
			return fmt.Errorf("failed to health check provider %s: %s", name, err)
		}

		ConfiguredProviders[name] = provider
	}

	return nil
}
