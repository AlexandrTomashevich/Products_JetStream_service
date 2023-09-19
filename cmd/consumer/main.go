package main

import (
	"database/sql"

	"encoding/json"
	"log"
	"os"
	"os/signal"
	"service/internal/models"
	"service/internal/repository"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://0.0.0.0:4222")
	if err != nil {
		panic("nats error: " + err.Error())
	}
	defer nc.Drain()

	js, _ := nc.JetStream()

	streamName := "EVENTS"

	js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	})

	log.Println("start subscription")

	_, err = js.AddConsumer(streamName, &nats.ConsumerConfig{
		Durable:        "event-processor",
		DeliverSubject: "products",
		DeliverGroup:   "event-processor",
		AckPolicy:      nats.AckExplicitPolicy,
	})
	if err != nil {
		log.Println("add consumer error", err)
		panic(err)
	}
	defer js.DeleteConsumer(streamName, "event-processor")

	dbConnectionString := "user=postgres dbname=postgreSQL sslmode=disable password=" + os.Getenv("DB_PASSWORD")
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Println("postgres connect error", err)
		panic(err)
	}

	repo := repository.NewRepository()

	log.Println("start consumer")
	var wg sync.WaitGroup

	sub, err := js.QueueSubscribe("events.products", "event-processor", func(m *nats.Msg) {
		wg.Add(1)
		defer wg.Done()

		if m == nil {
			return
		}

		log.Printf("get message: %+v\n", m)

		var p models.Product
		if err := json.Unmarshal(m.Data, &p); err != nil {
			log.Println("unmarshal error", err)
			return
		}
		if err := repo.AddProduct(db, p); err != nil {
			log.Println("add product error", err)
		}
		m.Ack()
	}, nats.ManualAck())

	if err != nil {
		panic("subscribe error: " + err.Error())
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh

	sub.Unsubscribe()
	wg.Wait()

	log.Println("stop consumer")
}
