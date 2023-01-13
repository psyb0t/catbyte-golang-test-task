package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

var rabbitMQURI = "amqp://user:password@localhost:7001/"
var rabbitMQConn *amqp.Connection

type message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type response struct {
	Message string `json:"message"`
}

func messageHandler(c *gin.Context) {
	var msg message
	var resp response

	if err := c.ShouldBindJSON(&msg); err != nil {
		resp.Message = err.Error()

		c.JSON(http.StatusBadRequest, resp)

		return
	}

	if msg.Sender == "" || msg.Receiver == "" || msg.Message == "" {
		resp.Message = "Fields should not be empty"

		c.JSON(http.StatusBadRequest, resp)

		return
	}

	if err := sendMessageToRabbitMQ(msg); err != nil {
		resp.Message = err.Error()

		c.JSON(http.StatusInternalServerError, resp)

		return
	}

	resp.Message = "OK"
	c.JSON(http.StatusOK, resp)
}

func sendMessageToRabbitMQ(msg message) error {
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

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	log.Println("publishing message to rabbitmq", string(b))
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		return err
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

	r := gin.Default()
	r.POST("/message", messageHandler)

	log.Fatal(r.Run())
}
