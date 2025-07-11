package main

import (
	"context"
	"fmt"
	"log"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"bitbucket.org/junglee_games/getsetgo/eventqueue/eventqueuefactory"
)

func main() {
	kafka, err := eventqueuefactory.GetPublisher(&configs.DefaultPublisherConfig{Name: "kafka", Kafka: configs.DefaultKafkaConfig{Brokers: "localhost:29092", Topic: "invoices_to_pdf", PublisherCount: 100, AsyncQueueSize: 1000}})
	if err != nil {
		log.Panic(err)
	}
	n := 1000
	ch := kafka.GetAsyncPublishResponseChan()
	for i := 0; i < n; i++ {
		kafka.PublishAsync(context.Background(), fmt.Sprintf("%d", i), fmt.Sprintf(" message %d", i))
	}
	for i := 0; i < n; i++ {
		err = <-ch
		if err != nil {
			log.Panic(err)
		}
		log.Print("success")
	}
	kafka.Close()
}
