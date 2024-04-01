package provider

import (
	"context"
	"fmt"
	"net"
	"time"

	huaweihg659 "github.com/axatol/go-huawei-hg659"
	"github.com/axatol/home-assistant-integrations/pkg/broker"
	"github.com/axatol/home-assistant-integrations/pkg/homeassistant"
	"github.com/axatol/home-assistant-integrations/pkg/util"
	"github.com/codingconcepts/env"
	"github.com/rs/zerolog/log"
)

type HuaweiHG659Provider struct {
	configured bool
	client     *huaweihg659.Client
	ticker     *time.Ticker

	Enabled    bool          `env:"HUAWEI_HG659_ENABLED" default:"false"`
	PollRate   time.Duration `env:"HUAWEI_HG659_POLL_RATE" default:"60s"`
	Address    string        `env:"HUAWEI_HG659_ADDRESS" required:"true"`
	EntityName string        `env:"HUAWEI_HG659_ENTITY_NAME" default:"huawei_hg659"`
}

func (p *HuaweiHG659Provider) Name() string {
	return "huawei_hg65"
}

func (p *HuaweiHG659Provider) Configure() error {
	if err := env.Set(p); err != nil {
		log.Warn().Err(fmt.Errorf("failed to set env values: %s", err)).Send()
	}

	if !p.Enabled {
		return nil
	}

	client, err := huaweihg659.NewClient(
		p.Address,
		huaweihg659.WithHTTPRoundTriper(&util.LogRoundTripper{Name: "huawei_hg659"}),
	)

	if err != nil {
		return fmt.Errorf("failed to create huawei hg659 client: %s", err)
	}

	p.ticker = time.NewTicker(p.PollRate)
	p.client = client

	return nil
}

func (p *HuaweiHG659Provider) StateTopic() string {
	return fmt.Sprintf("homeassistant_integrations/%s/state", p.EntityName)
}

func (p *HuaweiHG659Provider) AvailabilityTopic() string {
	return fmt.Sprintf("homeassistant_integrations/%s/availability", p.EntityName)
}

func (p *HuaweiHG659Provider) Schema() map[string]homeassistant.EntityConfiguration {
	if !p.Enabled {
		return nil
	}

	stateTopic := p.StateTopic()

	device := homeassistant.DeviceInformation{
		Name:         "Huawei HG659",
		Identifiers:  []string{"huawei_router_hg659"},
		Manufacturer: "Huawei",
		Model:        "HG659",
	}

	availability := []homeassistant.Availability{{
		Topic:               p.AvailabilityTopic(),
		PayloadAvailable:    util.Ptr("online"),
		PayloadNotAvailable: util.Ptr("offline"),
	}}

	return map[string]homeassistant.EntityConfiguration{
		fmt.Sprintf("homeassistant/binary_sensor/%s/router_internet_connected/config", p.EntityName): {
			Name:          "Internet Connected",
			UniqueID:      "router_internet_connected",
			ValueTemplate: util.Ptr("{{ value_json.internet_connected }}"),
			StateTopic:    stateTopic,
			DeviceClass:   util.Ptr("power"),
			Device:        device,
			Availability:  availability,
			Origin:        deviceOrigin,
		},
		fmt.Sprintf("homeassistant/sensor/%s/router_internet_self_test_message/config", p.EntityName): {
			Name:           "Self-test Message",
			UniqueID:       "router_internet_self_test_message",
			ValueTemplate:  util.Ptr("{{ value_json.self_test_message }}"),
			StateTopic:     stateTopic,
			EntityCategory: util.Ptr("diagnostic"),
			Device:         device,
			Availability:   availability,
			Origin:         deviceOrigin,
		},
		fmt.Sprintf("homeassistant/sensor/%s/router_internet_connection_status/config", p.EntityName): {
			Name:          "Internet Connection Status",
			UniqueID:      "router_internet_connection_status",
			ValueTemplate: util.Ptr("{{ value_json.internet_connection_status }}"),
			StateTopic:    stateTopic,
			Device:        device,
			Availability:  availability,
			Origin:        deviceOrigin,
		},
		fmt.Sprintf("homeassistant/sensor/%s/router_internet_err_reason/config", p.EntityName): {
			Name:           "Internet Err Reason",
			UniqueID:       "router_internet_err_reason",
			ValueTemplate:  util.Ptr("{{ value_json.internet_err_reason }}"),
			StateTopic:     stateTopic,
			Device:         device,
			EntityCategory: util.Ptr("diagnostic"),
			Availability:   availability,
			Origin:         deviceOrigin,
		},
		fmt.Sprintf("homeassistant/sensor/%s/router_internet_uptime/config", p.EntityName): {
			Name:              "Internet Uptime",
			UniqueID:          "router_internet_uptime",
			ValueTemplate:     util.Ptr("{{ value_json.internet_uptime }}"),
			StateTopic:        stateTopic,
			StateClass:        util.Ptr("total_increasing"),
			DeviceClass:       util.Ptr("duration"),
			UnitOfMeasurement: util.Ptr("ms"),
			Device:            device,
			Availability:      availability,
			Origin:            deviceOrigin,
		},
		fmt.Sprintf("homeassistant/sensor/%s/router_device_uptime/config", p.EntityName): {
			Name:              "Device Uptime",
			UniqueID:          "router_device_uptime",
			ValueTemplate:     util.Ptr("{{ value_json.device_uptime }}"),
			StateTopic:        stateTopic,
			StateClass:        util.Ptr("total_increasing"),
			DeviceClass:       util.Ptr("duration"),
			UnitOfMeasurement: util.Ptr("ms"),
			Device:            device,
			Availability:      availability,
			Origin:            deviceOrigin,
		},
	}
}

func (p *HuaweiHG659Provider) Health(ctx context.Context) (map[string]any, error) {
	info, err := p.client.GetDeviceInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get router device info: %s", err)
	}

	meta := map[string]any{"up_time": info.UpTime}

	return meta, nil
}

func (p *HuaweiHG659Provider) Subscribe(ctx context.Context) <-chan broker.Payload {
	sink := make(chan broker.Payload, PROVIDER_SINK_BUFFER_SIZE)
	go p.run(ctx, sink)
	return sink
}

func (p *HuaweiHG659Provider) run(ctx context.Context, sink chan<- broker.Payload) {
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
					Err(fmt.Errorf("failed to poll huawei hg659: %s", err)).
					Msg("error polling huawei hg659")

				sink <- broker.Payload{Topic: p.AvailabilityTopic(), Data: "offline"}
			} else {
				sink <- broker.Payload{Topic: p.AvailabilityTopic(), Data: "online"}
			}
		}
	}
}

func (p *HuaweiHG659Provider) poll(ctx context.Context, sink chan<- broker.Payload) error {
	info, err := p.client.GetDeviceInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get router device info: %s", err)
	}

	diagnosis, err := p.client.GetInternetDiagnosis(ctx)
	if err != nil {
		return fmt.Errorf("failed to get internet diagnosis: %s", err)
	}

	lookupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, lookupErr := net.DefaultResolver.LookupHost(lookupCtx, "www.tpg.com.au")

	connected := "OFF"
	if info.UpTime > 0 && diagnosis.ConnectionStatus == "Connected" && lookupErr == nil {
		connected = "ON"
	}

	message := "None"
	if lookupErr != nil {
		message = lookupErr.Error()
	}

	sink <- broker.Payload{
		Topic: p.StateTopic(),
		Data: map[string]any{
			"internet_connected":         connected,
			"self_test_message":          message,
			"internet_connection_status": diagnosis.ConnectionStatus,
			"internet_err_reason":        diagnosis.ErrReason,
			"internet_uptime":            diagnosis.Uptime,
			"device_uptime":              info.UpTime,
		},
	}

	return nil
}
