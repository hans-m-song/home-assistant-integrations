package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/axatol/go-zeversolar"
	"github.com/axatol/home-assistant-integrations/pkg/homeassistant"
	"github.com/axatol/home-assistant-integrations/pkg/util"
	"github.com/codingconcepts/env"
	"github.com/rs/zerolog/log"
)

type ZeverSolarTLC5000Provider struct {
	configured bool
	ticker     *time.Ticker
	client     *zeversolar.Client
	healthy    bool
	lastPoint  map[string]any

	Enabled       bool          `env:"ZEVER_SOLAR_TLC5000_ENABLED" default:"false"`
	PollRate      time.Duration `env:"ZEVER_SOLAR_TLC5000_POLL_RATE" default:"30s"`
	RouterAddress string        `env:"ZEVER_SOLAR_TLC5000_ADDRESS" required:"true"`
	EntityName    string        `env:"ZEVER_SOLAR_TLC5000_ENTITY_NAME" default:"zever_solar_tlc5000"`
}

func (p *ZeverSolarTLC5000Provider) Name() string {
	return "zever_solar_tlc5000"
}

func (p *ZeverSolarTLC5000Provider) Configure() error {
	if err := env.Set(p); err != nil {
		log.Warn().Err(fmt.Errorf("failed to set env values: %s", err)).Send()
	}

	if !p.Enabled {
		return nil
	}

	if p.RouterAddress == "" {
		return fmt.Errorf("router address is required")
	}

	p.ticker = time.NewTicker(p.PollRate)
	p.client = &zeversolar.Client{
		Address: p.RouterAddress,
		Client:  &http.Client{Transport: &util.LogRoundTripper{Name: "zeversolar"}},
	}

	return nil
}

func (p *ZeverSolarTLC5000Provider) StateTopic() string {
	return fmt.Sprintf("homeassistant_integrations/%s/state", p.EntityName)
}

func (p *ZeverSolarTLC5000Provider) AvailabilityTopic() string {
	return fmt.Sprintf("homeassistant_integrations/%s/availability", p.EntityName)
}

func (p *ZeverSolarTLC5000Provider) Interval() <-chan time.Time {
	if p.ticker == nil {
		return nil
	}

	return p.ticker.C
}

func (p *ZeverSolarTLC5000Provider) DeviceMetadata() homeassistant.DeviceInformation {
	return homeassistant.DeviceInformation{
		Identifiers:  []string{"zeversolar_inverter_tlc5000"},
		Name:         "Solar Inverter",
		Manufacturer: "Zeversolar",
		Model:        "TLC5000",
	}
}

func (p *ZeverSolarTLC5000Provider) EntityConfigurationSet() map[string]homeassistant.EntityConfiguration {
	builder := &homeassistant.EntityConfigurationSet{
		StateTopic:        p.StateTopic(),
		AvailabilityTopic: p.AvailabilityTopic(),
		Device:            p.DeviceMetadata(),
	}

	builder.Add("sensor", p.Name(), "solar_last_updated", homeassistant.EntityConfiguration{
		Name:           "Solar Last Updated",
		DeviceClass:    "timestamp",
		EntityCategory: "diagnostic",
	})

	builder.Add("sensor", p.Name(), "solar_power_ac_w", homeassistant.EntityConfiguration{
		Name:              "Solar Power AC (W)",
		StateClass:        "measurement",
		DeviceClass:       "power",
		UnitOfMeasurement: "W",
	})

	builder.Add("sensor", p.Name(), "solar_energy_today_kwh", homeassistant.EntityConfiguration{
		Name:              "Solar Energy Today (kWh)",
		StateClass:        "total_increasing",
		LastReset:         util.Midnight().Format(time.RFC3339),
		DeviceClass:       "energy",
		UnitOfMeasurement: "kWh",
	})

	builder.Add("sensor", p.Name(), "solar_status", homeassistant.EntityConfiguration{
		Name:        "Solar Status",
		DeviceClass: "power",
	})

	return builder.Entities()
}

func (p *ZeverSolarTLC5000Provider) Poll(ctx context.Context) (map[string]any, error) {
	point, err := p.client.GetInverterData(ctx)
	if err != nil {
		p.lastPoint = nil

		if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "host is down") {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get inverter data: %s", err)
	}

	status := "OFF"
	if point.Status == "OK" {
		status = "ON"
	}

	data := map[string]any{
		"solar_last_updated":     point.Timestamp.Format(time.RFC3339),
		"solar_power_ac_w":       point.PowerAC,
		"solar_energy_today_kwh": point.EnergyToday,
		"solar_status":           status,
	}

	p.lastPoint = data

	return data, nil
}

func (p *ZeverSolarTLC5000Provider) Health(ctx context.Context) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	point, err := p.client.GetInverterData(ctx)
	if err != nil {
		p.lastPoint = nil
		return nil, fmt.Errorf("failed to get inverter data: %s", err)
	}

	return map[string]any{"energy_today": point.EnergyToday}, nil
}
