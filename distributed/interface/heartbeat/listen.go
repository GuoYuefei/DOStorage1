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

// exclude 修复时可能会使用exclude, 防止重复选取data server
func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers()
	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidates = append(candidates, s)
		}
	}

	length := len(candidates)
	if length < n {
		return
	}
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	return
}

