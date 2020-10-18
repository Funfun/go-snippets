package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Address string
	Topic   string
	conn    *kafka.Conn
}

func NewConsumer(addr, topic string) (*Consumer, error) {
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", addr, topic, partition)
	if err != nil {
		return nil, fmt.Errorf("failed to dial leader: %s", err)
	}

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	return &Consumer{Address: addr, Topic: topic, conn: conn}, nil
}

func (c Consumer) Read() (messages []string, err error) {
	batch := c.conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	defer func() {
		if err = batch.Close(); err != nil {
			return
		}
	}()

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		messages = append(messages, string(b))
	}

	return
}

func (c Consumer) Close() error {
	return c.conn.Close()
}
