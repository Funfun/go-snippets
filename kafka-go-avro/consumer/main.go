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

	consumer, err := NewConsumer(*kafkaAddr, *kafkaTopic)
	if err != nil {
		log.Fatalf("failed to create consumer: %s", err)
	}

	r := gin.Default()
	r.GET("/messages", func(c *gin.Context) {
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
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
