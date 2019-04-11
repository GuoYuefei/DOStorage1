package main

import (
	"github.com/streadway/amqp"
	"log"
	"storage/rabbitMQ_study"
)

var r123 chan int = rabbitMQ_study.RandomInt123(2000,1800, 1, 2,4)

// 发布和订阅其实是一个一对多的关系
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	rabbitMQ_study.FailOnError(err, "connection is error")
	defer conn.Close()

	ch, err := conn.Channel()
	rabbitMQ_study.FailOnError(err, "注册通道失败")
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
	rabbitMQ_study.FailOnError(err, "declare exchange error")
	body := []byte(bodyproduct())

	err = ch.Publish(
		"logs",			//向logs发布消息
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: body,
		})
	rabbitMQ_study.FailOnError(err, "发布消息失败")

	log.Printf("[x] send %s", string(body))
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