package homeassistant

type Entity struct {
	Name          string `json:"name"`
	UniqueID      string `json:"unique_id"`
	ValueTemplate string `json:"value_template"`
	StateTopic    string `json:"state_topic"`
	DeviceClass   string `json:"device_class"`
	Device        Device `json:"device"`
}

type Device struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
}
