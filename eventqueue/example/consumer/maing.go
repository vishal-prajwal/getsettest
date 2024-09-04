package main

import (
	"context"
	"fmt"
	"log"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"bitbucket.org/junglee_games/getsetgo/eventqueue/eventqueuefactory"
)

func main() {

	c, err := eventqueuefactory.GetConsumer(&configs.DefaultConsumerConfig{
		Kafka: configs.DefaultKafkaConfig{
			Brokers:        "localhost:9092",
			Topic:          "jwr-image-masking-qa10",
			GroupId:        "group4",
			AsyncQueueSize: 1000,
			BatchSize:      2,
			MaxWaitSeconds: 10,
		},
		Name: "KAFKA",
	})
	if err != nil {
		log.Panic("Error in Making consumer ", err)
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
		msgs, err := c.ReadMessageWithUnCommit(context.Background())
		if err != nil {
			log.Panic(err)
		}
		// for _, m := range msgs {
		// 	log.Print(string(m.Key), string(m.Value))
		// }
		log.Printf("messages fetched %v", string(msgs.Key))

		err = c.Commit(context.Background())
		fmt.Println("Committing ", err)
		// time.Sleep(time.Second * 5)
	}
}
