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

func defaultPublishHandler(client mqtt.Client, msg mqtt.Message) {
	log.Trace().
		Str("topic", msg.Topic()).
		Uint16("message_id", msg.MessageID()).
		Bytes("payload", msg.Payload()).
		Msg("unhandled message")
}

type Broker struct {
	client mqtt.Client
}

type BrokerOption func(*mqtt.ClientOptions) *mqtt.ClientOptions

func WithWill(will Payload) (BrokerOption, error) {
	var (
		raw []byte
		ok  bool
		err error
	)

	if raw, ok = will.Data.([]byte); !ok {
		if raw, err = json.Marshal(will.Data); err != nil {
			return nil, err
		}
	}

	return func(co *mqtt.ClientOptions) *mqtt.ClientOptions {
		return co.SetBinaryWill(will.Topic, raw, 1, true)
	}, nil
}

// NewBroker creates a new MQTT broker.
//
// addr must be a fully qualified URI, e.g. tcp://localhost:1883.
//
// clientID will be cut to 24 characters
func NewBroker(addr string, options ...BrokerOption) (*Broker, error) {
	clientID := fmt.Sprintf("hass_int_%d", rand.Int64())[:24]

	opts := mqtt.NewClientOptions().
		AddBroker(addr).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetDefaultPublishHandler(defaultPublishHandler)

	for _, option := range options {
		opts = option(opts)
	}

	client := mqtt.NewClient(opts)
	if err := handleToken(client.Connect()); err != nil {
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
		if err := b.Publish(ctx, payload.Topic, payload.Data); err != nil {
			log.Error().Any("payload", payload).Err(err).Send()
		}
	}
}

type MQTTPublishOptions struct {
	QoS      byte
	Retained bool
}

func WithQoS(qos byte) func(*MQTTPublishOptions) {
	return func(o *MQTTPublishOptions) { o.QoS = qos }
}

func WithRetained(retained bool) func(*MQTTPublishOptions) {
	return func(o *MQTTPublishOptions) { o.Retained = retained }
}

func (b *Broker) Publish(ctx context.Context, topic string, payload any, options ...func(*MQTTPublishOptions)) error {
	opts := &MQTTPublishOptions{QoS: 1, Retained: false}
	for _, opt := range options {
		opt(opts)
	}

	var (
		raw []byte
		err error
	)

	log := log.Ctx(ctx).With().
		Str("topic", topic).
		Str("raw", string(raw)).
		Int("qos", int(opts.QoS)).
		Bool("retained", opts.Retained).
		Logger()

	if str, ok := payload.(string); ok {
		raw = []byte(str)
		log = log.With().Str("payload", string(raw)).Logger()
	} else {
		raw, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %s", err)
		}
		log = log.With().RawJSON("payload", raw).Logger()
	}

	// token := b.client.Publish(topic, opts.QoS, opts.Retained, raw)
	// if err := handleToken(token); err != nil {
	// 	log.Trace().Err(fmt.Errorf("failed to publish message: %s", err)).Send()
	// 	return fmt.Errorf("failed to publish to %s: %s", topic, err)
	// }

	log.Trace().Msg("published message")
	return nil
}

func (b *Broker) Subscribe(topic string, handler func(topic string, payload any)) error {
	callback := func(c mqtt.Client, m mqtt.Message) { handler(m.Topic(), m.Payload()) }
	token := b.client.Subscribe(topic, 1, callback)
	if err := handleToken(token); err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %s", topic, err)
	}

	return nil
}

func handleToken(token mqtt.Token) error {
	if !token.WaitTimeout(time.Second * 5) {
		return fmt.Errorf("timed out")
	}

	if token.Error() != nil {
		return token.Error()
	}

	return nil
}
