package heartbeat

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/rabbitmq"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string] time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(config.Pub.RABBITMQ_SERVER)
	defer q.Close()

	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, err := strconv.Unquote(string(msg.Body))
		utils.FailOnError(err, "Unquote error")
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {

	for {
		time.Sleep(6 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(12 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}

}

func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()

	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}

	return  ds
}

func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)

	if n == 0 {
		return ""
	}

	return ds[rand.Intn(n)]
}




