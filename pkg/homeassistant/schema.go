package homeassistant

type DeviceInformation struct {
	Name         string   `json:"name"`
	Identifiers  []string `json:"identifiers,omitempty"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
}

type Availability struct {
	Topic               string  `json:"topic"`
	ValueTemplate       *string `json:"value_template,omitempty"`
	PayloadAvailable    *string `json:"payload_available,omitempty"`
	PayloadNotAvailable *string `json:"payload_not_available,omitempty"`
}

// https://www.home-assistant.io/integrations/sensor.mqtt/
type EntityConfiguration struct {
	Name              string            `json:"name"`
	UniqueID          string            `json:"unique_id"`
	ValueTemplate     *string           `json:"value_template,omitempty"`
	StateTopic        string            `json:"state_topic"`
	StateClass        *string           `json:"state_class,omitempty"`
	LastReset         *string           `json:"last_reset,omitempty"`
	DeviceClass       *string           `json:"device_class,omitempty"`
	UnitOfMeasurement *string           `json:"unit_of_measurement,omitempty"`
	EntityCategory    *string           `json:"entity_category,omitempty"`
	Origin            EntityOrigin      `json:"origin,omitempty"`
	Device            DeviceInformation `json:"device"`
	AvailabilityMode  *string           `json:"availability_mode,omitempty"`
	Availability      []Availability    `json:"availability,omitempty"`
}

type EntityOrigin struct {
	Name            string `json:"name"`
	SoftwareVersion string `json:"sw_version"`
	SupportURL      string `json:"support_url"`
}
