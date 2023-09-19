package main

import (
	"encoding/json"
	"log"
	"service/internal/models"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://0.0.0.0:4222")
	if err != nil {
		panic(err)
	}
	defer nc.Drain()

	js, err := nc.JetStream()
	if err != nil {
		log.Println("create jetstream error", err)
		panic(err)
	}

	streamName := "EVENTS"

	js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	})

	var testProducts = []models.Product{
		{
			Name:         "pen",
			Price:        10,
			Manufacturer: "bic",
		},
		{
			Name:         "toy",
			Price:        100,
			Manufacturer: "lego",
		},
	}

	log.Println("producer start")

	for _, p := range testProducts {
		body, _ := json.Marshal(p)
		_, err := js.Publish("events.products", body)
		if err != nil {
			log.Println("publish error", err)
		}
	}

	log.Println("producer end")
}
