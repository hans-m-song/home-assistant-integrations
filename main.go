package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/axatol/home-assistant-integrations/pkg/broker"
	"github.com/axatol/home-assistant-integrations/pkg/config"
	"github.com/axatol/home-assistant-integrations/pkg/homeassistant"
	"github.com/axatol/home-assistant-integrations/pkg/provider"
	"github.com/axatol/home-assistant-integrations/pkg/server"
	"github.com/rs/zerolog/log"
)

func init() {
	config.Configure()

	log.Debug().
		Any("config", config.Values).
		Send()

	log.Info().
		Str("build_commit", config.BuildCommit).
		Str("build_time", config.BuildTime).
		Str("build_version", config.BuildVersion).
		Msg("starting")
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, os.Kill)

	log.Debug().Msg("configuring mqtt broker")
	mqtt, err := broker.NewBroker(config.Values.MQTTURI)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("failed to create broker: %s", err)).Send()
	}

	log.Debug().Msg("configuring providers")
	if err := provider.Configure(ctx); err != nil {
		log.Fatal().Err(fmt.Errorf("failed to configure providers: %s", err)).Send()
	}

	log.Debug().Msg("starting providers")
	if err := provider.Start(ctx, mqtt); err != nil {
		log.Fatal().Err(fmt.Errorf("failed to start providers: %s", err)).Send()
	}

	log.Debug().Msg("configuring server")
	mux := server.Configure()
	server := http.Server{
		Handler: mux.Router,
		Addr:    fmt.Sprintf(":%d", config.Values.ListenPort),
	}

	for _, provider := range provider.Providers {
		log.Debug().Str("provider", provider.Name()).Msg("adding provider as probe target")
		mux.AddProbeTarget(provider)
	}

	log.Info().Msgf("listening on %s", server.Addr)
	go server.ListenAndServe()

	select {
	case s := <-sig:
		log.Info().Str("signal", s.String()).Msg("received signal, shutting down gracefully")
		cancel(nil)
	case <-ctx.Done():
		if err := context.Cause(ctx); err != nil && err != context.Canceled {
			log.Error().Err(fmt.Errorf("shutting down gracefully: %s", err)).Send()
		} else {
			log.Info().Msg("shutting down gracefully")
		}
	}

	if err := mqtt.Publish(ctx, homeassistant.BRIDGE_AVAILABILITY_TOPIC, "offline"); err != nil {
		log.Error().Err(fmt.Errorf("failed to set bridge availability to offline: %s", err)).Send()
	}

	ctx, cancel = context.WithCancelCause(context.Background())
	defer cancel(nil)

	go func() {
		if err := server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			cancel(fmt.Errorf("failed to shut down server gracefully: %s", err))
		} else {
			cancel(nil)
		}
	}()

	select {
	case <-time.After(time.Second * 5):
		log.Error().Err(fmt.Errorf("failed to shut down gracefully: timed out")).Send()
		os.Exit(1)
	case s := <-sig:
		log.Error().Str("signal", s.String()).Err(fmt.Errorf("failed to shut down gracefully: received signal")).Send()
		os.Exit(1)
	case <-ctx.Done():
		if err := context.Cause(ctx); err != nil && err != context.Canceled {
			log.Error().Err(fmt.Errorf("failed to shut down gracefully: %s", err)).Send()
			os.Exit(1)
		}
	}
}
