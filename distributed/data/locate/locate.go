package locate

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/rabbitmq"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, e := os.Stat(name)
	return !os.IsNotExist(e)
}

func StartLocate() {
	q := rabbitmq.New(config.Pub.RABBITMQ_SERVER)
	defer q.Close()

	q.Bind("dataServers")
	c := q.Consume()

	for msg := range c {
		object, err := strconv.Unquote(string(msg.Body))
		utils.FailOnError(err, "Unquote error")

		if Locate(config.ServerData.STORAGE_ROOT+"/objects/"+object) {
			q.Send(msg.ReplyTo, config.ServerData.LISTEN_ADDRESS)
		}
	}
}