package main

import (
	"context"
	"log"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"bitbucket.org/junglee_games/getsetgo/eventqueue/eventqueuefactory"
)

func main() {

	c, err := eventqueuefactory.GetConsumer(&configs.DefaultConsumerConfig{
		Kafka: configs.DefaultKafkaConfig{
			Brokers:        "localhost:29092;localhost:39092",
			Topic:          "invoices_to_pdf",
			GroupId:        "group4",
			AsyncQueueSize: 1000,
			BatchSize:      50000,
			MaxWaitSeconds: 10,
		},
	})
	if err != nil {
		log.Panic(err)
	}
	// for i := 0; i < 10; i++ {
	// 	log.Print("reading message")
	// 	msg, err := c.ReadMessage(context.Background())
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	log.Print(string(msg.Key), string(msg.Value))
	// }
	for {
		log.Print("reading batch")
		msgs, err := c.ReadBatch(context.Background())
		if err != nil {
			log.Panic(err)
		}
		// for _, m := range msgs {
		// 	log.Print(string(m.Key), string(m.Value))
		// }
		log.Printf("%d messages fetched", len(msgs))
		// time.Sleep(time.Second * 5)
	}
}
