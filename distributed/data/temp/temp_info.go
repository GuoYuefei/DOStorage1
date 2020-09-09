package temp

import (
	"encoding/json"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"os"
	"path"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string					// is hash
	Size int64
}

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func (t *tempInfo) writeToFile() error {
	f, e := os.Create(path.Join(config.TempRoot, t.Uuid))
	if e != nil {
		return e
	}
	defer f.Close()
	b, _ := json.Marshal(t)
	f.Write(b)
	return nil
}