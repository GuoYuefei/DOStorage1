package main

import (
	"github.com/streadway/amqp"
	"log"
	"storage/rabbitMQ_study"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	rabbitMQ_study.FailOnError(err, "连接失败")
	defer conn.Close()

	ch, err := conn.Channel()
	rabbitMQ_study.FailOnError(err, "fail to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	rabbitMQ_study.FailOnError(err, "fail to declare an exchange")

	queue, err := ch.QueueDeclare(
		"",
		true,
		false,
		true,
		false,
		nil,
	)
	rabbitMQ_study.FailOnError(err, "fail to declare a queue")

	err = ch.QueueBind(
		queue.Name,
		"",
		"logs",
		false,
		nil,
	)
	rabbitMQ_study.FailOnError(err, "fail to bind a queue")

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	rabbitMQ_study.FailOnError(err, "fali to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("[x] receive %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}
