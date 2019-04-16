package heartbeat

import (
	"os"
	"storage/distributed/rabbitmq"
	"time"
)

func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	ticker := time.NewTicker(5 * time.Second)

	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		<- ticker.C
	}

}
