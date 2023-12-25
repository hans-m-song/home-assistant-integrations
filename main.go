package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/axatol/home-assistant-integrations/pkg/config"
	"github.com/axatol/home-assistant-integrations/pkg/provider"
	"github.com/axatol/home-assistant-integrations/pkg/server"
	"github.com/rs/zerolog/log"
)

func init() {
	config.Configure()
	log.Debug().Fields(config.Values()).Send()
}

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		cancel(fmt.Errorf("received signal %s", sig))
	}()

	mux := server.Configure()
	server := http.Server{Handler: mux.Router, Addr: fmt.Sprintf(":%d", config.ListenPort)}
	go func() {
		log.Info().Msgf("listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cancel(fmt.Errorf("server closed unexpectedly: %s", err))
		}
	}()

	log.Trace().Msg("configuring providers")
	if err := provider.Configure(ctx, config.ProviderConfig); err != nil {
		cancel(fmt.Errorf("failed to configure providers: %s", err))
	}

	log.Trace().Msg("adding configured providers to probe targets")
	for _, v := range provider.ConfiguredProviders {
		mux.AddProbeTarget(v)
	}

	select {
	case sig := <-sigs:
		log.Error().Err(fmt.Errorf("received signal %s", sig)).Send()
	case <-ctx.Done():
		if err := context.Cause(ctx); err != nil && err != context.Canceled {
			log.Error().Err(err).Msg("shutting down gracefully")
		} else {
			log.Info().Err(err).Msg("shutting down gracefully")
		}
	}

	ctx, cancel = context.WithCancelCause(context.Background())
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Trace().Msg("shutting down server")
		if err := server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			log.Error().Err(fmt.Errorf("failed to shut down server: %s", err)).Send()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for name, provider := range provider.ConfiguredProviders {
			log.Trace().Str("provider", name).Msg("shutting down provider")
			if err := provider.Close(); err != nil {
				log.Error().Err(fmt.Errorf("failed to shut down provider: %s", err)).Send()
			}
		}
	}()

	go func() {
		wg.Wait()
		log.Trace().Msg("graceful shutdown complete")
		cancel(nil)
	}()

	select {
	case <-time.After(time.Second * 5):
		log.Error().Err(fmt.Errorf("failed to shut down gracefully: timed out")).Send()
		os.Exit(1)
	case sig := <-sigs:
		log.Error().Err(fmt.Errorf("failed to shut down gracefully: received signal %s", sig)).Send()
		os.Exit(1)
	case <-ctx.Done():
		if err := context.Cause(ctx); err != nil && err != context.Canceled {
			log.Error().Err(fmt.Errorf("failed to shut down gracefully: %s", err)).Send()
			os.Exit(1)
		}
	}
}
