package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	UpdatedAt string `json:"updated_at"`
}

func main() {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{"localhost:9093"}, config)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("user_changes", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Listening for user changes...")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var user User
			fmt.Println("Received message:", string(msg.Value))
			if err := json.Unmarshal(msg.Value, &user); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			fmt.Printf("Received change: ID=%d, Name=%s, Email=%s, UpdatedAt=%s\n",
				user.ID, user.Name, user.Email, user.UpdatedAt)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)
		}
	}
}
