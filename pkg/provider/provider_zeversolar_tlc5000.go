package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/axatol/go-zeversolar"
	"github.com/axatol/home-assistant-integrations/pkg/broker"
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

func (p *ZeverSolarTLC5000Provider) Schema() map[string]homeassistant.EntityConfiguration {
	if !p.Enabled {
		return nil
	}

	stateTopic := p.StateTopic()

	device := homeassistant.DeviceInformation{
		Identifiers:  []string{"zeversolar_inverter_tlc5000"},
		Name:         "Solar Inverter",
		Manufacturer: "Zeversolar",
		Model:        "TLC5000",
	}

	availability := []homeassistant.Availability{{
		Topic:               p.AvailabilityTopic(),
		PayloadAvailable:    util.Ptr("online"),
		PayloadNotAvailable: util.Ptr("offline"),
	}}

	return map[string]homeassistant.EntityConfiguration{
		fmt.Sprintf("homeassistant/sensor/%s/solar_last_updated/config", p.EntityName): {
			Name:           "Solar Last Updated",
			UniqueID:       "solar_last_updated",
			ValueTemplate:  util.Ptr("{{ value_json.last_updated }}"),
			StateTopic:     stateTopic,
			DeviceClass:    util.Ptr("timestamp"),
			Device:         device,
			EntityCategory: util.Ptr("diagnostic"),
			Availability:   availability,
			Origin:         deviceOrigin,
		},
		fmt.Sprintf("homeassistant/sensor/%s/solar_power_ac_w/config", p.EntityName): {
			Name:              "Solar Power AC (W)",
			UniqueID:          "solar_power_ac_w",
			ValueTemplate:     util.Ptr("{{ value_json.power_ac }}"),
			StateTopic:        stateTopic,
			StateClass:        util.Ptr("measurement"),
			DeviceClass:       util.Ptr("power"),
			UnitOfMeasurement: util.Ptr("W"),
			Device:            device,
			Availability:      availability,
			Origin:            deviceOrigin,
		},
		fmt.Sprintf("homeassistant/sensor/%s/solar_energy_today_kwh/config", p.EntityName): {
			Name:              "Solar Energy Today (kWh)",
			UniqueID:          "solar_energy_today_kwh",
			ValueTemplate:     util.Ptr("{{ value_json.energy_today }}"),
			StateTopic:        stateTopic,
			StateClass:        util.Ptr("total_increasing"),
			LastReset:         util.Ptr(util.Midnight().Format(time.RFC3339)),
			DeviceClass:       util.Ptr("energy"),
			UnitOfMeasurement: util.Ptr("kWh"),
			Device:            device,
			Availability:      availability,
			Origin:            deviceOrigin,
		},
		fmt.Sprintf("homeassistant/binary_sensor/%s/solar_status/config", p.EntityName): {
			Name:          "Solar Status",
			UniqueID:      "solar_status",
			ValueTemplate: util.Ptr("{{ value_json.status }}"),
			StateTopic:    stateTopic,
			DeviceClass:   util.Ptr("power"),
			Device:        device,
			Availability:  availability,
			Origin:        deviceOrigin,
		},
	}
}

func (p *ZeverSolarTLC5000Provider) Health(ctx context.Context) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sink := make(chan broker.Payload, 1)
	defer close(sink)

	if err := p.poll(ctx, sink); err != nil {
		return nil, fmt.Errorf("failed to get inverter data: %s", err)
	}

	item := <-sink
	data, ok := item.Data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("failed to cast data to map[string]any")
	}

	return data, nil
}

func (p *ZeverSolarTLC5000Provider) Subscribe(ctx context.Context) <-chan broker.Payload {
	sink := make(chan broker.Payload, PROVIDER_SINK_BUFFER_SIZE)
	go p.run(ctx, sink)
	return sink
}

func (p *ZeverSolarTLC5000Provider) run(ctx context.Context, sink chan<- broker.Payload) {
	for {
		select {
		case <-ctx.Done():
			p.ticker.Stop()
			sink <- broker.Payload{Topic: p.AvailabilityTopic(), Data: "offline"}
			close(sink)
			return

		case <-p.ticker.C:
			if err := p.poll(ctx, sink); err != nil {
				log.Error().
					Err(fmt.Errorf("failed to poll zever solar tlc5000: %s", err)).
					Msg("error polling zever solar tlc5000")

				sink <- broker.Payload{Topic: p.AvailabilityTopic(), Data: "offline"}
			} else {
				sink <- broker.Payload{Topic: p.AvailabilityTopic(), Data: "online"}
			}
		}
	}
}

func (p *ZeverSolarTLC5000Provider) poll(ctx context.Context, sink chan<- broker.Payload) error {
	point, err := p.client.GetInverterData(ctx)
	if err != nil {
		p.lastPoint = nil
		return fmt.Errorf("failed to get inverter data: %s", err)
	}

	status := "OFF"
	if point.Status == "OK" {
		status = "ON"
	}

	p.lastPoint = map[string]any{
		"last_updated": point.Timestamp.Format(time.RFC3339),
		"power_ac":     point.PowerAC,
		"energy_today": point.EnergyToday,
		"status":       status,
	}

	sink <- broker.Payload{Topic: p.StateTopic(), Data: p.lastPoint}

	return nil
}
