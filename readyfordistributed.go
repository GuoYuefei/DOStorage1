package main

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	//node rabbitmq 的前置准备， 声明两个exchange
	readyForRabbitMQ()

	// elasticsearch 的前置准备， 创建metadata索引以及object类型的映射
	// node elasticsearch 7.x不在支持指定索引类型
	readyForElastic()
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

func readyForRabbitMQ() {
	conn, _ := amqp.Dial(os.Getenv("RABBITMQ_SERVER"))
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()

	declare(ch, "dataServers")
	declare(ch, "apiServers")
}

func readyForElastic() {
	client := http.Client{}
	//b := strings.NewReader(`
	//		{
	//			"mappings":{
	//				"objects":{
	//					"properties":{
	//						"name":{"type":"string","index":"not_analyzed"},
	//						"version":{"type":"integer"},
	//						"size":{"type":"integer"},
	//						"hash":{"type":"string"}
	//					}
	//				}
	//			}
	//		}
	//	`)

	b := strings.NewReader(`
		{
			"mappings":{
				"properties":{
					"name":{"type":"text","index": true},
					"version":{"type":"integer"},
					"size":{"type":"integer"},
					"hash":{"type":"text"}
				}
			}
		}
	`)

	r, e := http.NewRequest(http.MethodPut, "http://"+config.Pub.ES_SERVER+"/metadata", b)
	if e != nil {
		log.Fatal(e)
	}
	r.Header.Set("Content-Type", "application/json")
	response, e := client.Do(r)
	if e != nil {
		log.Fatal(e)
	}
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		log.Fatalf("fail to create metadata index, the error code is %s,\n%s", response.Status, string(body))
	}

	log.Println("ok!")

}
