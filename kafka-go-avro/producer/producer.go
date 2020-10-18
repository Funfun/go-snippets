package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Address string
	Topic   string
	conn    *kafka.Conn
}

func NewProducer(addr, topic string) (*Producer, error) {
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", addr, topic, partition)
	if err != nil {
		return nil, fmt.Errorf("failed to dial leader: %s", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	return &Producer{Address: addr, Topic: topic, conn: conn}, nil
}

func (p Producer) Send(msg []byte) error {
	_, err := p.conn.WriteMessages(
		kafka.Message{Value: []byte(msg)},
	)
	if err != nil {
		return fmt.Errorf("failed to write messages: %s", err)
	}

	return nil
}

func (p Producer) Close() error {
	return p.conn.Close()
}
