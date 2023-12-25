package broker

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

const (
	clientIDMaxLength = 23
)

func init() {
	mqtt.ERROR = mqttLogger{"ERROR"}
	mqtt.CRITICAL = mqttLogger{"CRITICAL"}
	mqtt.WARN = mqttLogger{"WARN"}
	mqtt.DEBUG = mqttLogger{"DEBUG"}
}

var (
	_ mqtt.Logger = (*mqttLogger)(nil)
)

type mqttLogger struct{ level string }

func (l mqttLogger) Println(v ...interface{}) {
	log := log.Trace().Str("namespace", "mqtt").Str("level", l.level)
	for _, item := range v {
		log = log.Interface("item", item)
	}
	log.Send()
}

func (l mqttLogger) Printf(format string, v ...interface{}) {
	log.Trace().Str("namespace", "mqtt").Str("level", l.level).Msgf(format, v...)
}

func defaultMQTTPublishHandler(client mqtt.Client, msg mqtt.Message) {
	log.Trace().
		Str("namespace", "mqtt").
		Str("topic", msg.Topic()).
		Uint16("message_id", msg.MessageID()).
		Bytes("payload", msg.Payload()).
		Msg("unhandled message")
}

type LastWill struct {
	Topic   string
	Payload any
}

type mqttBrokerImpl struct {
	clientID string
	client   mqtt.Client
}

func generateClientID(purpose string) string {
	clientID := strings.Builder{}
	clientID.WriteString("hass_")
	if len(purpose) > 16 {
		purpose = purpose[:16]
	}
	clientID.WriteString(purpose)
	clientID.WriteString("_")
	id := make([]byte, clientIDMaxLength-clientID.Len())
	if _, err := rand.Read(id); err != nil {
		log.Warn().Err(fmt.Errorf("failed to generate unique portion of client id: %s", err)).Send()
	}
	clientID.Write(id)
	return clientID.String()
}

// NewBroker creates a new MQTT broker.
//
// addr must be a fully qualified URI, e.g. tcp://localhost:1883.
//
// purpose will be cut to 16 characters to generate a client id in the form
// `hai_<purpose>_<at least 3 random bytes>`.
func NewMQTTBroker(addr, purpose string, will *LastWill) (Broker, error) {
	clientID := generateClientID(purpose)
	opts := mqtt.NewClientOptions().
		AddBroker(addr).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetDefaultPublishHandler(defaultMQTTPublishHandler)

	if will != nil {
		raw, err := json.Marshal(will.Payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal will payload: %s", err)
		}

		opts.SetBinaryWill(will.Topic, raw, 1, true)
	}

	client := mqtt.NewClient(opts)
	if err := handleMQTTToken(client.Connect()); err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err)
	}

	return &mqttBrokerImpl{clientID, client}, nil
}

func (b *mqttBrokerImpl) Name() string {
	return fmt.Sprintf("broker:mqtt:%s", b.clientID)
}

func (b *mqttBrokerImpl) Health(ctx context.Context) (map[string]any, error) {
	if !b.client.IsConnected() {
		return nil, fmt.Errorf("not connected")
	}

	return map[string]any{"connected": true}, nil
}

func (b *mqttBrokerImpl) Disconnect(deadline time.Duration) error {
	b.client.Disconnect(uint(deadline.Milliseconds()))
	return nil
}

func (b *mqttBrokerImpl) Publish(topic string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %s", err)
	}

	if err := handleMQTTToken(b.client.Publish(topic, 1, false, raw)); err != nil {
		return fmt.Errorf("failed to publish to %s: %s", topic, err)
	}

	return nil
}

func (b *mqttBrokerImpl) Subscribe(topic string, handler func(topic string, payload any)) error {
	if err := handleMQTTToken(b.client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) { handler(m.Topic(), m.Payload()) })); err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %s", topic, err)
	}

	return nil
}

func handleMQTTToken(token mqtt.Token) error {
	if !token.WaitTimeout(time.Second * 5) {
		return fmt.Errorf("timed out")
	}

	if token.Error() != nil {
		return token.Error()
	}

	return nil
}
