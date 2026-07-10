package devices

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"open-fermentations/internal/env"

	mqttLib "github.com/eclipse/paho.mqtt.golang"
)

type Payload struct {
	Hello string `json:"hello"`
}

type MQTT interface {
	Publish(topic string, payload *Payload) error
	Subscribe(topic string, qos byte, messageHandler mqttLib.MessageHandler) error
}

type mqtt struct {
	env    *env.Env
	client mqttLib.Client
}

// Subscribe implements [MQTT].
func (m mqtt) Subscribe(topic string, qos byte, messageHandler mqttLib.MessageHandler) error {
	if m.client.IsConnected() == false {
		return fmt.Errorf("mqtt client is not connected")
	}

	token := m.client.Subscribe(topic, qos, messageHandler)

	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}

	slog.Info("Successfully subscribed to topic", slog.String("topic", topic), slog.Int("qos", int(qos)))
	return nil
}

// Publish implements [MQTT].
func (m mqtt) Publish(topic string, payload *Payload) error {
	byteData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("marshalling payload for publish", slog.Group("mqtt", slog.String("topic", topic)))
		return err
	}
	token := m.client.Publish(topic, byte(1), false, string(byteData))
	token.Wait()
	return nil
}

var mqttInstance *mqtt

func New(env *env.Env) *mqtt {
	if mqttInstance == nil {
		broker := fmt.Sprintf("mqtt://%s:%s", env.Mqtt.Host, env.Mqtt.Port)
		opts := mqttLib.NewClientOptions().
			AddBroker(broker).
			SetClientID(env.Mqtt.ClientID).
			SetUsername(env.Mqtt.User).
			SetPassword(env.Mqtt.Password).
			SetCleanSession(false)
		client := mqttLib.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			slog.Error("Could not connect to mqtt on initial startup", slog.Group("mqtt",
				slog.String("broker", broker),
				slog.String("username", env.Mqtt.User),
				slog.String("clientId", env.Mqtt.ClientID),
			), slog.String("error", token.Error().Error()))
		}

		mqttInstance = &mqtt{
			env,
			client,
		}
	}

	return mqttInstance
}

var _ MQTT = mqtt{}
