package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"storage/rabbitMQ_study"
)

func main() {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	rabbitMQ_study.FailOnError(err, "Fail to connect to RabbitMQ server")
	defer conn.Close()

	ch, err := conn.Channel()
	rabbitMQ_study.FailOnError(err, "Fail to open a Channel")
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

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true, //独家享用
		false,
		nil,
	)
	rabbitMQ_study.FailOnError(err, "Fail to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	for index, s := range os.Args[1:] {
		log.Printf("%d: Binding queue %s to exchange %s with routing key %s",
			index, queue.Name, "logs_direct", s)

		// use for to bind
		err = ch.QueueBind(
			queue.Name,
			s,
			"logs_direct",
			false,
			nil,
		)

		rabbitMQ_study.FailOnError(err, "Fail to bind")
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	rabbitMQ_study.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()


	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever


}



