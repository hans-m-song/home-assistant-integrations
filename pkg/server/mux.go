package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type ProbeTarget interface {
	Name() string
	Health(ctx context.Context) (map[string]any, error)
}

type Mux struct {
	mu           sync.RWMutex
	probeTargets []ProbeTarget

	Router *chi.Mux
}

func Configure(probeTargets ...ProbeTarget) *Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Compress(5, "application/json"))
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestLogger(requestLogFormatter{}))
	router.Use(middleware.Recoverer)

	server := Mux{
		Router:       router,
		probeTargets: probeTargets,
	}

	router.Group(func(r chi.Router) {
		r.Use(middlewareContentType("application/json"))
		r.Get("/readyz", server.readinessProbe)
	})

	return &server
}

func (server *Mux) AddProbeTarget(target ProbeTarget) *Mux {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.probeTargets = append(server.probeTargets, target)
	return server
}

type metadataItem struct {
	Healthy bool `json:"healthy"`
	Data    any  `json:"data"`
}

func (server *Mux) readinessProbe(w http.ResponseWriter, r *http.Request) {
	metadata := map[string]metadataItem{}
	errors := map[string]error{}

	wg := sync.WaitGroup{}
	server.mu.RLock()
	for _, target := range server.probeTargets {
		wg.Add(1)
		target := target
		go func() {
			defer wg.Done()
			name := target.Name()
			data, err := target.Health(r.Context())
			if err != nil {
				errors[name] = err
			}

			metadata[name] = metadataItem{
				Healthy: err == nil,
				Data:    data,
			}
		}()
	}
	wg.Wait()
	server.mu.RUnlock()

	raw, err := json.Marshal(metadata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Any("metadata", metadata).Msg("failed to marshal metadata")
		return
	}

	if len(errors) > 0 {
		log.Error().Any("errors", errors).Any("metadata", metadata).Send()
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(raw); err != nil {
		log.Error().Err(fmt.Errorf("failed to write response: %s", err)).Send()
	}
}
