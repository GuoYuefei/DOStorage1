package main

import (
	"github.com/GuoYuefei/DOStorage1/rabbitMQ_study"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	rabbitMQ_study.FailOnError(err, "Fail to connect to rabbitmq server")
	defer conn.Close()

	ch, err := conn.Channel()
	rabbitMQ_study.FailOnError(err, "Fail to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		false,
		true,
		false,
		false,
		nil,
	)
	rabbitMQ_study.FailOnError(err, "Fail to declare an exchange")

	body := rabbitMQ_study.GenerateMess(10)

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}



	err = ch.Publish(
		"logs_direct",
		os.Args[1],
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
	})

	rabbitMQ_study.FailOnError(err, "Fail to publish a message")

	log.Fatalf("[*] sent %s", string(body))

}
