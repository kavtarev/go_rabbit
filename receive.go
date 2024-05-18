package main

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main3() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Println(1111111, err)
	}

	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(222222, err)
	}

	err = channel.ExchangeDeclare("rrr-exchange", "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Println(222.22, err)
	}

	queue, err := channel.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		fmt.Println(333333, err)
	}

	err = channel.QueueBind(queue.Name, "", "rrr-exchange", false, nil)
	if err != nil {
		fmt.Println(3333.3555, err)
	}

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		fmt.Println(44444, err)
	}

	i := 0
	for msg := range messages {
		time.Sleep(2000)
		fmt.Println(i, string(msg.Body))
		i++
		err := msg.Ack(false)
		if err != nil {
			fmt.Println(434535)
		}
	}
}
