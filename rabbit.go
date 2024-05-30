package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	default_exchange = "default_rabbit_exchange"
	default_queue    = "default_queue"
)

func NewChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Println("error connected to rabbit")
	}

	channel, err := conn.Channel()
	if err != nil {
		fmt.Println("error get channel")
	}

	err = channel.ExchangeDeclare(default_exchange, "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Println("cannot declare"+default_exchange, err)
	}

	_, err = channel.QueueDeclare(default_queue, false, false, false, false, nil)
	if err != nil {
		fmt.Println(333333, err)
	}

	return channel
}
