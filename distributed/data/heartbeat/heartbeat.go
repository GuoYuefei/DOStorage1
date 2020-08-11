package heartbeat

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/rabbitmq"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(config.Pub.RABBITMQ_SERVER)
	defer q.Close()

	ticker := time.NewTicker(5 * time.Second)

	for {
		q.Publish("apiServers", config.ServerData.LISTEN_ADDRESS)
		<- ticker.C
	}

}
