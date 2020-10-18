package main

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Address string
	Topic   string
}

func NewConsumer(addr, topic string) (*Consumer, error) {
	return &Consumer{Address: addr, Topic: topic}, nil
}

func (c Consumer) Read() ([]string, error) {
	messages := []string{}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{c.Address},
		Topic:     c.Topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		messages = append(messages, string(m.Value))
	}

	return messages, r.Close()
}
