package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var (
	kafkaTopic = flag.String("topic", "default", "Kafka topic")
	kafkaAddr  = flag.String("addr", "localhost:9092", "Kafka connection address")
)

func newMessageHandler(c *gin.Context) {
	producer, err := NewProducer(*kafkaAddr, *kafkaTopic)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})

		return
	}
	defer producer.Close()

	msg := c.PostForm("message")

	if err := producer.Send([]byte(msg)); err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"status": "send",
		"info": gin.H{
			"topic": producer.Topic,
			"addr":  producer.Address,
		},
	})
}

func main() {
	flag.Parse()

	r := gin.Default()
	r.POST("/messages", newMessageHandler)
	r.Run()
}
