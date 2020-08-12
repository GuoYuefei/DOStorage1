package objectstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server	string
	Uuid	string
}

func NewTempPutStream(server, hash string, size int64) (*TempPutStream, error) {
	request, e := http.NewRequest(http.MethodPost, "http://"+server+"/temp/"+hash, nil)
	if e != nil {
		return nil, e
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	uuid, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return nil, e
	}
	return &TempPutStream{server, string(uuid)}, nil
}

// 以uuid访问数据节点的临时文件
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, e := http.NewRequest(http.MethodPatch, "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if e != nil {
		return 0, e
	}
	client := http.Client{}
	r, e := client.Do(request)
	if e != nil {
		return 0, e
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := http.MethodDelete
	if good {
		method = http.MethodPut
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}

