package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var (
	kafkaTopic = flag.String("topic", "default", "Kafka topic")
	kafkaAddr  = flag.String("addr", "localhost:9092", "Kafka connection address")
)

func messagesHandler(c *gin.Context) {
	consumer, err := NewConsumer(*kafkaAddr, *kafkaTopic)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})

		return
	}
	defer consumer.Close()

	m, err := consumer.Read()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"messages": m,
	})
}

func main() {
	flag.Parse()

	r := gin.Default()
	r.GET("/messages", messagesHandler)

	r.Run()
}
