package main

import (
	"log/slog"
	"open-fermentations/internal/devices"
	"open-fermentations/internal/env"

	"github.com/joho/godotenv"
)

const topic = "hello"

func main() {
	godotenv.Load()

	env := env.GetEnv()

	mq := devices.New(env)
	slog.Info("Publisher initialized and ready to publish on topic 'hello'.")

	slog.Info("Attempting to publish message.")
	if err := mq.Publish(topic, &devices.Payload{Hello: "world"}); err != nil {
		panic(err)
	}

	slog.Info("Successfully published message once")
}
