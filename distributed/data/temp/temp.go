package temp

import (
	"encoding/json"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/data/locate"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
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

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	utils.Log.Printf(utils.Debug, "will call function temp.%s", m)
	switch m {
	case http.MethodPut: put(w, r)
	case http.MethodPatch: patch(w, r)
	case http.MethodPost: post(w, r)
	case http.MethodDelete: delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	uu := uuid.New().String()
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{
		uu,
		name,
		size,
	}
	e = t.writeToFile()
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, e = os.Create(path.Join(config.TempRoot, t.Uuid+".dat"))
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(uu))
}

func patch(w http.ResponseWriter, r *http.Request) {
	uu := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uu)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := path.Join(config.TempRoot, uu)
	datFile := infoFile+".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	//f, e := os.Create(datFile)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, e = io.Copy(f, r.Body)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		utils.Log.Printf(utils.Debug, "actual size %d, exceeds %d", actual, tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	uu := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uu)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := path.Join(config.TempRoot, uu)
	datFile := infoFile+".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		utils.Log.Println(utils.Err, e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempinfo.Size {
		os.Remove(datFile)
		utils.Log.Printf(utils.Info, "actual size %d, exceeds %d", actual, tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(datFile, tempinfo)
	w.WriteHeader(http.StatusAccepted)
}

func delete(_ http.ResponseWriter, r *http.Request) {
	uu := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := path.Join(config.TempRoot, uu)
	datFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(datFile)
}

func readFromFile(uu string) (*tempInfo, error) {
	f, e := os.Open(path.Join(config.TempRoot, uu))
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}

// 最后的名字定位成 <hash>.X.<hash of shard X>
// todo 错误处理
func commitTempObject(datFile string, info *tempInfo) {
	f, e := os.Open(datFile)
	if e != nil {
		utils.Log.Println(utils.Exception, e)
	}
	d := url.PathEscape(utils.CalculateHash(f))
	e = f.Close()
	if e != nil {
		time.Sleep(50*time.Millisecond)
		f.Close()
	}
	// 最多花费1秒
	for i := 10; i > 0; i-- {
		e = os.Rename(datFile, path.Join(config.ObjectRoot, info.Name+"."+d))
		if e != nil {
			time.Sleep(100*time.Millisecond)
		} else {
			break
		}
	}

	locate.Add(info.hash(), info.id())
}
