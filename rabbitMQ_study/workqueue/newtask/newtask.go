package main

import (
	"github.com/GuoYuefei/DOStorage1/rabbitMQ_study"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var r123 chan int = rabbitMQ_study.RandomInt123(2000,1800, 1, 2,4)

// 生产任务
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrors(err, "连接失败")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnErrors(err, "Fail to open a channel")
	defer ch.Close()

	q,err := ch.QueueDeclare(
		"tasks",
		true,
		false,
		false,
		false,
		nil,
		)
	failOnErrors(err, "failed to declare a queue")
	for {
		body := bodyproduct()
		if body=="" {
			break
		}

		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType: "text/plain",
				Body: []byte(body),
			})
		failOnErrors(err, "Failed to publish message")
		log.Printf("[x] Sent %s", body)

		time.Sleep(1 * time.Second)
	}

}

func failOnErrors(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyproduct() string {
	total, open := <-r123
	s := ""
	if !open {
		return s
	}
	for i := 0; i < total; i++ {
		s += "."
	}

	return s
}