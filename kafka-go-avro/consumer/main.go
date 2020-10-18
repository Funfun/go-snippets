package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaTopic = flag.String("topic", "default", "Kafka topic")
	kafkaAddr  = flag.String("addr", "localhost:9092", "Kafka connection address")
)

func messagesHandler(c *gin.Context) {
	// to consume messages
	topic := c.Query("topic")
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{*kafkaAddr},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(42)

	var ms []string

	for {
		m, err := r.ReadMessage(c.Request.Context())
		if err != nil {
			break
		}
		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		ms = append(ms, string(m.Value))
	}

	if err := r.Close(); err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"messages": ms,
	})
}

func main() {
	flag.Parse()

	r := gin.Default()
	r.GET("/messages", messagesHandler)

	r.Run()
}
