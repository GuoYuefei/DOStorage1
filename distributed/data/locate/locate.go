package locate

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/rabbitmq"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var objects map[string]int = make(map[string]int)
var mutex sync.Mutex

var ObjectRoot = filepath.Join(config.ServerData.STORAGE_ROOT, "objects")
var TempRoot = filepath.Join(config.ServerData.STORAGE_ROOT, "temp")

func init() {
	// make dir for ./data/objects
	_, e := os.Stat(ObjectRoot)
	if e != nil {
		e := os.MkdirAll(ObjectRoot, os.ModePerm)
		if e != nil {
			log.Fatal(e)
		}
	}

	_, e = os.Stat(TempRoot)
	if e != nil {
		e := os.MkdirAll(TempRoot, os.ModePerm)
		if e != nil {
			log.Fatal(e)
		}
	}
}

func Locate(hash string) bool {
	mutex.Lock()
	_, ok := objects[hash]
	mutex.Unlock()
	return ok
}

func Add(hash string) {
	mutex.Lock()
	objects[hash] = 1
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
		if Locate(hash) {
			// 定位成功就发送自己的信息出去
			q.Send(msg.ReplyTo, config.ServerData.LISTEN_ADDRESS)
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(filepath.Join(ObjectRoot, "/*"))
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}