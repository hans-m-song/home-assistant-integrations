package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
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
	log := log.Trace().Str("namespace", "mqtt").Str("lvl", l.level)
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

type Broker struct {
	client mqtt.Client
}

// NewBroker creates a new MQTT broker.
//
// addr must be a fully qualified URI, e.g. tcp://localhost:1883.
//
// clientID will be cut to 24 characters
func NewMQTTBroker(addr string, will *LastWill) (*Broker, error) {
	clientID := fmt.Sprintf("hass_int_%d", rand.Int64())[:24]

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

	return &Broker{client}, nil
}

func (b *Broker) Disconnect(deadline time.Duration) error {
	b.client.Disconnect(uint(deadline.Milliseconds()))
	return nil
}

type Payload struct {
	Topic string
	Data  any
}

func (b *Broker) Listen(ctx context.Context, source <-chan Payload) {
	for payload := range source {
		if err := b.Publish(payload.Topic, payload.Data); err != nil {
			log.Error().Any("payload", payload).Msg("failed to publish")
		}
	}
}

type MQTTPublishOptions struct {
	QoS      byte
	Retained bool
}

func WithMQTTQoS(qos byte) func(*MQTTPublishOptions) {
	return func(o *MQTTPublishOptions) { o.QoS = qos }
}

func WithMQTTRetained(retained bool) func(*MQTTPublishOptions) {
	return func(o *MQTTPublishOptions) { o.Retained = retained }
}

func (b *Broker) Publish(topic string, payload any, options ...func(*MQTTPublishOptions)) error {
	opts := &MQTTPublishOptions{QoS: 1, Retained: false}
	for _, opt := range options {
		opt(opts)
	}

	var (
		raw []byte
		err error
	)

	if str, ok := payload.(string); ok {
		raw = []byte(str)
	} else {
		raw, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %s", err)
		}
	}

	log.Trace().
		Str("namespace", "mqtt").
		Str("topic", topic).
		Str("payload", string(raw)).
		Msg("publishing message")

	token := b.client.Publish(topic, opts.QoS, opts.Retained, raw)
	if err := handleMQTTToken(token); err != nil {
		return fmt.Errorf("failed to publish to %s: %s", topic, err)
	}

	return nil
}

func (b *Broker) Subscribe(topic string, handler func(topic string, payload any)) error {
	callback := func(c mqtt.Client, m mqtt.Message) { handler(m.Topic(), m.Payload()) }
	token := b.client.Subscribe(topic, 1, callback)
	if err := handleMQTTToken(token); err != nil {
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
