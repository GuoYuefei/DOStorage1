package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"time"
)

//var ra = rabbitMQ_study.RandomInt123(100000, 100000, 1, 3, 4)
var gorNum int = 1
var forever = make(chan bool)
var workerNameNext int = 0

// 工作的人
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "连接失败")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "开启通道失败")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks",
		true,
		false,
		false,
		false,
		nil,
		)
	failOnError(err, "声明队列失败")

	err = ch.Qos(
		1,
		0,
		false,
		)
	failOnError(err, "file to set Qos")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,	//手动回应
		false,
		false,
		false,
		nil,
		)
	failOnError(err, "Failed to register a consumer")

	workerNameNext++
	go work(msgs, workerNameNext)


	log.Printf("[*] Waiting for message. To exit press CTRL+C")
	<-forever

}



func work(msgs <-chan amqp.Delivery, workerName int) {
	var ti time.Duration
	isbreak := false
	taskNum := 1


	log.Printf("worker [#%d] is runing!\n", workerName)
	log.Printf("-------------------there are %d Goroutines", gorNum)
	begin := time.Now()
	for d := range msgs {
		ti = time.Since(begin)
		log.Println("********Time*******接收消息时间： ", ti)



		log.Printf("#%d Received a message: %s", workerName, d.Body)
		dot_count := bytes.Count(d.Body, []byte("."))
		//r := <- ra
		//if dot_count+r >= 5 {
		//	gorNum++
		//	go work(msgs, workerName+1, taskNum+1)
		//} else if dot_count+r <= 3 && taskNum > 4 {
		//	isbreak = true			// 如果出现任务不繁忙了就在让工人自杀（不人道）
		//}

		t := time.Duration(dot_count)

		if ti > (28+t*4)*time.Second {
			isbreak = true
		} else if (ti < 5*time.Second && taskNum<3) || ti < 18*time.Second {
			gorNum++
			workerNameNext++
			go work(msgs, workerNameNext)
		}

		time.Sleep(t * time.Second)
		d.Ack(false)
		log.Printf("worker [#%d]'s task %d Done!", workerName, taskNum)
		taskNum++
		if isbreak {
			gorNum--
			log.Printf("worker [#%d] is dying!", workerName)
			break
		}
		begin = time.Now()
	}
	log.Printf("----------------there are %d Goroutines", gorNum)
	if gorNum==0 {
		// 如果没协程了，那么主函数可以退出了
		forever <- false
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}