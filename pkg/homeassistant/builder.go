package homeassistant

import (
	"fmt"

	"github.com/axatol/home-assistant-integrations/pkg/config"
)

const (
	BRIDGE_AVAILABILITY_TOPIC = "homeassistant_integrations/bridge/availability"
	BRIDGE_NAME               = "Home Assistant Integrations"
	BRIDGE_SUPPORT_URL        = "https://github.com/axatol/home-assistant-integrations/issues"
)

type EntityConfigurationSet struct {
	StateTopic        string
	AvailabilityTopic string
	Device            DeviceInformation
	entities          map[string]EntityConfiguration
}

func (s *EntityConfigurationSet) Add(kind, name, id string, entity EntityConfiguration) *EntityConfigurationSet {
	topic := fmt.Sprintf("homeassistant/%s/%s/%s/config", kind, name, id)

	if entity.UniqueID == "" {
		entity.UniqueID = id
	}

	if entity.ValueTemplate == "" {
		entity.ValueTemplate = fmt.Sprintf("{{ value_json.%s }}", entity.UniqueID)
	}

	if entity.StateTopic == "" {
		entity.StateTopic = s.StateTopic
	}

	if entity.Availability == nil {
		entity.Availability = []Availability{
			{Topic: s.AvailabilityTopic},
			{Topic: BRIDGE_AVAILABILITY_TOPIC},
		}
	}

	if entity.AvailabilityMode == "" {
		entity.AvailabilityMode = "latest"
	}

	if entity.Device == nil {
		entity.Device = &s.Device
	}

	if entity.Origin == nil {
		entity.Origin = &EntityOrigin{
			Name:            BRIDGE_NAME,
			SupportURL:      BRIDGE_SUPPORT_URL,
			SoftwareVersion: config.BuildVersion,
		}
	}

	if entity.StateClass == "" {
		entity.StateClass = "measurement"
	}

	if s.entities == nil {
		s.entities = map[string]EntityConfiguration{}
	}

	s.entities[topic] = entity
	return s
}

func (s *EntityConfigurationSet) Entities() map[string]EntityConfiguration {
	return s.entities
}
