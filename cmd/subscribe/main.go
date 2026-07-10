package main

import (
	"encoding/json"
	"log/slog"
	"open-fermentations/internal/devices"
	"open-fermentations/internal/env"

	// Added import for time package
	mqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

const topic = "hello"

func main() {
	godotenv.Load()

	env := env.GetEnv()

	mq := devices.New(env)

	messageHandler := func(client mqttLib.Client, msg mqttLib.Message) {
		slog.Info("Message received")
		var payload devices.Payload
		err := json.Unmarshal(msg.Payload(), &payload)
		if err != nil {
			slog.Error("Failed to unmarshal payload", slog.String("error", err.Error()))
			return
		}
		slog.Info("Received message",
			slog.String("topic", msg.Topic()),
			slog.Any("payload", payload),
		)
	}

	qos := byte(1)

	if err := mq.Subscribe(topic, qos, messageHandler); err != nil {
		slog.Error("Failed to subscribe", slog.String("error", err.Error()))
		return
	}

	slog.Info("Subscription setup complete. Waiting for messages on topic:", slog.String("topic", topic))

	select {}
}
