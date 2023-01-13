package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	redisAddr   = "localhost:6379"
	redisClient *redis.Client
)

type message struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}

type response struct {
	Message string `json:"message"`
}

func getMessages(sender, receiver string) ([]message, error) {
	messages := make([]message, 0)

	length, err := redisClient.LLen(context.TODO(), "messages").Result()
	if err != nil {
		return nil, err
	}

	ms := make([]string, 0)
	for i := length - 1; i >= 0; i-- {
		m, err := redisClient.LIndex(context.TODO(), "messages", int64(i)).Result()
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	for _, m := range ms {
		var msg message
		if err := json.Unmarshal([]byte(m), &msg); err != nil {
			return nil, err
		}

		if msg.Sender == sender && msg.Receiver == receiver {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

func messageListHandler(c *gin.Context) {
	var resp response
	sender := c.Query("sender")
	receiver := c.Query("receiver")

	if sender == "" || receiver == "" {
		resp.Message = "sender or receiver params should not be empty"

		c.JSON(http.StatusBadRequest, resp)

		return
	}

	messages, err := getMessages(sender, receiver)
	if err != nil {
		resp.Message = err.Error()

		c.JSON(http.StatusInternalServerError, resp)

		return
	}

	c.JSON(http.StatusOK, messages)
}

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	r := gin.Default()
	r.GET("/message/list", messageListHandler)

	log.Fatal(r.Run(":8081"))
}
