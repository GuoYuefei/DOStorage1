package locate

import (
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/rabbitmq"
	"github.com/GuoYuefei/DOStorage1/distributed/types"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects map[string]int = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	q := rabbitmq.New(config.Pub.RABBITMQ_SERVER)
	defer q.Close()

	q.Bind("dataServers")
	c := q.Consume()

	for msg := range c {
		hash, err := strconv.Unquote(string(msg.Body))
		utils.PanicOnError(err, "Unquote error")
		utils.Log.Printf(utils.Debug, "get hash %s for locate\n", hash)
		id := Locate(hash)
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: config.ServerData.LISTEN_ADDRESS, Id: id})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(filepath.Join(config.ObjectRoot, "/*"))
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			utils.PanicOnError(fmt.Errorf("file name error, format error"), "")
			continue
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			utils.PanicOnError(fmt.Errorf("file name error. id is not a number"), "")
		}
		objects[hash] = id
	}
}
