package es

import (
	"encoding/json"
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var es_server string

func init() {
	es_server = config.Pub.ES_SERVER
}

type Metadata struct {
	Name    string		`json:"name"`
	Version int			`json:"version"`
	Size    int64		`json:"size"`
	Hash    string		`json:"hash"`
}

type hit struct {
	Source Metadata		`json:"_source"`
}

type total struct {
	Value int			`json:"value"`
	Relation string		`json:"relation"`
}

type hits struct {
	Total total			`json:"total"`
	Hits []hit			`json:"hits"`
}

type searchResult struct {
	Hits hits			`json:"hits"`
}

func getMetadata(name string, versionId int) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d/_source", es_server, name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return
	}
	json.Unmarshal(result, &meta)
	return
}

func SearchLatestVersion(name string) (meta Metadata, e error) {
	uri := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&" +
		"size=1&sort=version:desc", es_server, url.PathEscape(name))
	r, e := http.Get(uri)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create",
		es_server, name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	request.Header.Set("Content-Type", "application/json")
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

// from, size 指定分页和显示个数
// 若name 为空，则查找所有对象的版本
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		es_server, from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}