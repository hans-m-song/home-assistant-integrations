package provider

import (
	"context"
	"fmt"

	advantageair "github.com/axatol/go-advantage-air"
	"github.com/axatol/home-assistant-integrations/pkg/broker"
	"github.com/axatol/home-assistant-integrations/pkg/util"
)

type AdvantageAirHubProvider struct {
	aa   advantageair.Client
	mqtt broker.Broker
}

func (p *AdvantageAirHubProvider) Name() string {
	return "advantage_air_hub_provider"
}

func (p *AdvantageAirHubProvider) Configure(options map[string]any) error {
	hubAddress, err := util.GetByKey[string](options, "hub_address")
	if err != nil {
		return fmt.Errorf("hub_address is required: %s", err)
	}

	p.aa = advantageair.NewClient(hubAddress, 3)

	mqttAddress, err := util.GetByKey[string](options, "mqtt_address")
	if err != nil {
		return fmt.Errorf("mqtt_address is required: %s", err)
	}

	p.mqtt, err = broker.NewMQTTBroker(mqttAddress, "advantaair", &broker.LastWill{})
	if err != nil {
		return fmt.Errorf("failed to create mqtt broker: %s", err)
	}

	return nil
}

func (p *AdvantageAirHubProvider) Health(ctx context.Context) (map[string]any, error) {
	info, err := p.aa.GetSystemInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get hub system info: %s", err)
	}

	return map[string]any{"info": info}, nil
}

func (p *AdvantageAirHubProvider) Close() error {
	return p.mqtt.Disconnect(1000)
}
