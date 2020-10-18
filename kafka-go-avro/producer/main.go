package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	kafkaTopic = flag.String("topic", "default", "Kafka topic")
	kafkaAddr  = flag.String("addr", "localhost:9092", "Kafka connection address")
)

func main() {
	flag.Parse()

	producer, err := NewProducer(*kafkaAddr, *kafkaTopic)
	if err != nil {
		log.Fatalf("failed to create producer: %s", err)
	}
	defer producer.Close()

	r := gin.Default()
	r.POST("/messages", func(c *gin.Context) {
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
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
