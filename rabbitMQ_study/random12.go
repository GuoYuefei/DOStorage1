package rabbitMQ_study

import (
	"log"
	"time"
)

func RandomInt123(times int, seconds time.Duration, num1, num2, num3 int) (chan int) {
	timeout := make(chan bool)

	go func() {
		time.Sleep(seconds * time.Second)
		timeout <- true
	}()
	ch := make(chan int, 1024)
	c := 0

	go func() {
		for ; c < times; {
			select {
			case ch <- num1:
			case ch <- num2:
			case ch <- num3:

			}
			c++
		}
		timeout <- true
	}()

	go func() {
		// 如果timeout不阻塞了，就可以关闭channel了
		<-timeout
		close(ch)
	}()

	return ch

}


func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}