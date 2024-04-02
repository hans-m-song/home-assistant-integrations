package provider

import (
	"context"
	"fmt"
	"net"
	"time"

	huaweihg659 "github.com/axatol/go-huawei-hg659"
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

func (p *HuaweiHG659Provider) Interval() <-chan time.Time {
	if p.ticker == nil {
		return nil
	}

	return p.ticker.C
}

func (p *HuaweiHG659Provider) DeviceMetadata() homeassistant.DeviceInformation {
	return homeassistant.DeviceInformation{
		Name:         "Huawei HG659",
		Identifiers:  []string{"huawei_router_hg659"},
		Manufacturer: "Huawei",
		Model:        "HG659",
	}
}

func (p *HuaweiHG659Provider) EntityConfigurationSet() map[string]homeassistant.EntityConfiguration {
	builder := homeassistant.EntityConfigurationSet{
		AvailabilityTopic: p.AvailabilityTopic(),
		StateTopic:        p.StateTopic(),
		Device:            p.DeviceMetadata(),
	}

	builder.Add("binary_sensor", p.Name(), "router_internet_connected", homeassistant.EntityConfiguration{
		Name:        "Internet Connected",
		DeviceClass: "power",
	})

	builder.Add("sensor", p.Name(), "router_internet_self_test_message", homeassistant.EntityConfiguration{
		Name:           "Self-test Message",
		EntityCategory: "diagnostic",
	})

	builder.Add("sensor", p.Name(), "router_internet_connection_status", homeassistant.EntityConfiguration{
		Name: "Internet Connection Status",
	})

	builder.Add("sensor", p.Name(), "router_internet_err_reason", homeassistant.EntityConfiguration{
		Name:           "Internet Err Reason",
		EntityCategory: "diagnostic",
	})

	builder.Add("sensor", p.Name(), "router_internet_uptime", homeassistant.EntityConfiguration{
		Name:              "Internet Uptime",
		StateClass:        "total_increasing",
		DeviceClass:       "duration",
		UnitOfMeasurement: "ms",
	})

	builder.Add("sensor", p.Name(), "router_device_uptime", homeassistant.EntityConfiguration{
		Name:              "Device Uptime",
		StateClass:        "total_increasing",
		DeviceClass:       "duration",
		UnitOfMeasurement: "ms",
	})

	return builder.Entities()
}

func (p *HuaweiHG659Provider) Poll(ctx context.Context) (map[string]any, error) {
	info, err := p.client.GetDeviceInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get router device info: %s", err)
	}

	diagnosis, err := p.client.GetInternetDiagnosis(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get internet diagnosis: %s", err)
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

	data := map[string]any{
		"router_internet_connected":         connected,
		"router_internet_self_test_message": message,
		"router_internet_connection_status": diagnosis.ConnectionStatus,
		"router_internet_err_reason":        diagnosis.ErrReason,
		"router_internet_uptime":            diagnosis.Uptime,
		"router_device_uptime":              info.UpTime,
	}

	return data, nil
}

func (p *HuaweiHG659Provider) Health(ctx context.Context) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	info, err := p.client.GetDeviceInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get router device info: %s", err)
	}

	return map[string]any{"up_time": info.UpTime}, nil
}
