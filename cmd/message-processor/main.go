package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

var (
	rabbitMQURI  = "amqp://user:password@localhost:7001/"
	rabbitMQConn *amqp.Connection

	redisAddr   = "localhost:6379"
	redisClient *redis.Client
)

func consumeRabbitMQQueue() error {
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"messages",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		if err := redisClient.RPush(context.TODO(), "messages", msg.Body).Err(); err != nil {
			log.Printf("Failed to append message to Redis: %v", err)

			continue
		}

		log.Printf("Received message: %s", string(msg.Body))
	}

	return nil
}

func main() {
	var err error
	rabbitMQConn, err = amqp.Dial(rabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQConn.Close()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	log.Println("consuming rabbit mq queue")
	log.Fatal(consumeRabbitMQQueue())
}
