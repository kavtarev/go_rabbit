package main

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	default_exchange = "default_rabbit_exchange"
	default_queue    = "default_queue"
	exchange_uniq    = "exchange_uniq"
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
		fmt.Println("cant declare queue", err)
	}

	return channel
}

func Receiver(index int) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		fmt.Println("error connected to rabbit receiver")
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		fmt.Println("error get channel receiver")
	}
	defer channel.Close()

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Println("error Qos receiver")
	}

	ch, err := channel.Consume(default_queue, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("error cant consume receiver")
	}

	fmt.Printf("in goroutine %d\n", index)
	fin := make(chan struct{})

	go func() {
		if index == 2 {
			time.Sleep(time.Duration(10) * time.Second)
			fin <- struct{}{}
		}
	}()

	for {
		select {
		case msg := <-ch:
			time.Sleep(time.Duration(index) * time.Second)
			fmt.Printf("receive in %d %v\n", index, string(msg.Body))
			//if index == 1 {
			msg.Ack(false)
			// 	fmt.Printf("ack in %d\n", index)
			// }
		case <-fin:
			fmt.Printf("in channel %d, should close", index)
			return
		}
	}
}

func ReceiverWithExchange() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	// fanout разветвление
	err = ch.ExchangeDeclare(exchange_uniq, "fanout", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare("", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(q.Name, "", exchange_uniq, false, nil)
	if err != nil {
		panic(err)
	}

	msg, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for s := range msg {
		fmt.Println(string(s.Body))
	}

}
