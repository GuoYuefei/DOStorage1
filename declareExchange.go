package main

import "github.com/streadway/amqp"

func main() {
	conn, _ := amqp.Dial("amqp://test:test@localhost:5672/")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()

	declare(ch, "dataServers")
	declare(ch, "apiServers")

}

func declare(ch *amqp.Channel, name string)  {
	// rabbitmqctl list_exchanges 查看结果
	_ = ch.ExchangeDeclare(
		name,
		"fanout",
		true,
		false,
		false,
		false,
		nil)
}
