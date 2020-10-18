package main

import (
	"context"
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaTopic = flag.String("topic", "my-topic", "Kafka topic")
	kafkaAddr  = flag.String("addr", "localhost:9092", "Kafka connection address")
)

func newMessageHandler(c *gin.Context) {
	topic := c.PostForm("topic")
	producer, err := NewProducer(*kafkaAddr, topic)
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

func createTopic(c *gin.Context) {
	// to create topics
	topic := c.PostForm("topic")
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", *kafkaAddr, topic, partition)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"status": "created",
		"info": gin.H{
			"topic": topic,
			"addr":  *kafkaAddr,
		},
	})
}

func listTopics(c *gin.Context) {
	conn, err := kafka.Dial("tcp", *kafkaAddr)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		c.JSON(500, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}

	c.JSON(200, gin.H{
		"info": gin.H{
			"topics": m,
			"addr":   *kafkaAddr,
		},
	})
}

func main() {
	flag.Parse()

	r := gin.Default()
	r.POST("/messages", newMessageHandler)
	r.POST("/topic", createTopic)
	r.GET("/topics", listTopics)
	r.Run(":8082")
}
