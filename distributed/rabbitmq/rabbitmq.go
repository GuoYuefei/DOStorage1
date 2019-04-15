package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"storage/distributed/doslog"
)

type RabbitMQ struct {
	channel *amqp.Channel
	Name string
	exchange string
}

func New(s string) *RabbitMQ {
	conn, err := amqp.Dial(s)
	doslog.FailOnError(err, "Fail to connect to rabbitMQ server")

	ch, err := conn.Channel()
	doslog.FailOnError(err, "Fail to open a Channel")

	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	doslog.FailOnError(err, "Fail to declare a queue")

	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name

	return mq

}

func (q *RabbitMQ) Bind(exchange string) {
	err := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	doslog.FailOnError(err, "Fail to bind a queue")

	q.exchange = exchange
}

func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	doslog.FailOnError(err, "Fail to serialize body")

	err = q.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body: []byte(str),
		})

	doslog.FailOnError(err, "send message error")

}

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	doslog.FailOnError(err, "Fail to serialize body")

	err = q.channel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		})
	doslog.FailOnError(err, "Fail to publish")
}

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	c, err := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	doslog.FailOnError(err, "Fail to get message")

	return c
}

//Closer
func (q *RabbitMQ) Close() {
	_ = q.channel.Close()
}




