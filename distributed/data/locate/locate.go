package locate

import (
	"github.com/GuoYuefei/DOStorage1/distributed/doslog"
	"github.com/GuoYuefei/DOStorage1/distributed/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, e := os.Stat(name)
	return !os.IsNotExist(e)
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	q.Bind("dataServers")
	c := q.Consume()

	for msg := range c {
		object, err := strconv.Unquote(string(msg.Body))
		doslog.FailOnError(err, "Unquote error")

		if Locate(os.Getenv("STORAGE_ROOT")+"/objects/"+object) {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}