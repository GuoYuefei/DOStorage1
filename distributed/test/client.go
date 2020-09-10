package test

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)
const filepath = "./"
var apihost = "http://" + config.ServerInf.LISTEN_ADDRESS + "/objects/"
var api = "http://" + config.ServerInf.LISTEN_ADDRESS
var tokenloc = "token_location.tmp"

func Put(file string, ok bool) error {
	if ok {
		return putCorrect(file)
	}
	return putIncorrect(file)
}

func Get(name string, out io.Writer) error {
	resp, err := http.Get(apihost + name)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response error, %s\n", resp.Status)
	}
	if out == nil {
		out = os.Stdout
	}
	io.Copy(out, resp.Body)
	fmt.Println()
	return nil
}

func GetAllVersion(name string) error {
	resp, err := http.Get("http://" + config.ServerInf.LISTEN_ADDRESS + "/versions/" + name)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response error, %s\n", resp.Status)
	}
	io.Copy(os.Stdout, resp.Body)
	fmt.Println()
	return nil
}

func Del(name string) error {
	request, err := http.NewRequest(http.MethodDelete, apihost + name, nil)
	if err != nil {
		return err
	}
	client := http.Client{}
	do, err := client.Do(request)
	if err != nil {
		return err
	}
	if do.StatusCode != http.StatusOK {
		return fmt.Errorf("response error, %s\n", do.Status)
	}
	io.Copy(os.Stdout, do.Body)
	fmt.Println()
	return nil
}

func Locat(name string) error {
	do, err := http.Get("http://" + config.ServerInf.LISTEN_ADDRESS + "/locate/" + name)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if do.StatusCode != http.StatusOK {
		return fmt.Errorf("response error, %s\n", do.Status)
	}
	io.Copy(os.Stdout, do.Body)
	fmt.Println()
	return nil
}

func putCorrect(file string) error {
	r, e := http.NewRequest(http.MethodPut, apihost + file, nil)
	if e != nil {
		return e
	}
	f, e := os.Open(filepath+file)
	if e != nil {
		return e
	}
	defer f.Close()

	hash := sha256.New()
	//piper, pipew := io.Pipe()				// Pipe() 创建的管道是同步的，不适合这里的情况
	written, e := io.Copy(hash, f)
	if e != nil {
		return e
	}
	_, e = f.Seek(0, 0)
	if e != nil {
		return e
	}
	fmt.Println(written)

	sum := hash.Sum(nil)
	sha := base64.StdEncoding.EncodeToString(sum)
	utils.Log.Println(utils.Info, "sha ", sha)
	r.Body = f
	r.ContentLength = written
	r.Header.Set("Digest", "SHA-256="+sha)
	//r.Header.Set("Content-Length", ...)		// 无效

	client := http.Client{}
	response, e := client.Do(r)
	if e != nil {
		return e
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response status fault, message is %s", response.Status)
	}
	io.Copy(os.Stdout, response.Body)
	fmt.Println()
	return nil
}

func putIncorrect(file string) error {
	r, e := http.NewRequest(http.MethodPut, apihost + file, nil)
	if e != nil {
		return e
	}
	f, e := os.Open(filepath+file)
	if e != nil {
		return e
	}
	uu := uuid.New().String()
	r.Header.Set("Digest", "SHA-256="+uu)
	r.Body = f
	sta, e := f.Stat()
	if e != nil {
		return e
	}
	r.ContentLength = sta.Size()

	client := http.Client{}
	response, e := client.Do(r)
	if e != nil {
		return e
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response status fault, message is %s", response.Status)
	}
	io.Copy(os.Stdout, response.Body)
	fmt.Println()
	return nil
}

// size 第一次上传的大小
// 断点续传测试函数
func putBigFile(file string, size int64) error {
	utils.Log.Printf(utils.Info, "open file, %s\n", file)
	f, e := os.Open(filepath + file)
	if e != nil {
		return e
	}
	defer f.Close()
	utils.Log.Println(utils.Info, "open file suss!")
	hash := sha256.New()
	written, e := io.Copy(hash, f)
	if e != nil {
		return e
	}
	_, e = f.Seek(0, 0)		// 将指针移动到文件头
	if e != nil {
		return e
	}
	sum := hash.Sum(nil)
	sha := base64.StdEncoding.EncodeToString(sum)
	utils.Log.Println(utils.Info, "sha ", sha)
	client := http.Client{}
	//first post info
	utils.Log.Println(utils.Info, "first post info...")
	r, e := http.NewRequest(http.MethodPost, apihost+file, nil)
	r.Header.Set("Digest", "SHA-256="+sha)
	r.Header.Set("Size", strconv.FormatInt(written, 10))

	response, e := client.Do(r)
	if e != nil {
		return e
	}
	if response.StatusCode == http.StatusOK {
		return nil
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("%s", response.Status)
	}
	// post over
	location := response.Header.Get("location")
	utils.Log.Printf(utils.Info, "get upload uri: \n%s\n\tpost over\n", location)
	locationF, e := os.Create(tokenloc)
	if e != nil {
		return e
	}
	locationF.Write([]byte(location))
	locationF.Close()
	r, e = http.NewRequest(http.MethodPut, api+location, nil)
	if size > written {
		size = written
	}
	sectionR := io.NewSectionReader(f, 0, size)
	if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
		return e
	}
	r.Body = ioutil.NopCloser(sectionR)
	//r.Header.Set("Range", "bytes=0-"+size)		// 第一块可不加
	utils.Log.Printf(utils.Info, "do first put, size is %d\n", size)
	r.ContentLength = size
	res, e := client.Do(r)
	if e != nil {
		return e
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", res.Status)
	}
	// 第一次do put over
	utils.Log.Printf(utils.Info, "first put over!\n")
	time.Sleep(2*time.Second)

	locationb, _ := ioutil.ReadFile(tokenloc)
	location1 := string(locationb)
	//一次 head 看真实服务器端已经存储了多少
	reqH, _ := http.NewRequest(http.MethodHead, api+location1, nil)
	do, e := client.Do(reqH)
	if e != nil {
		return e
	}
	cl, e := strconv.ParseInt(do.Header.Get("Content-Length"), 10, 64)
	if e != nil {
		return e
	}
	utils.Log.Printf(utils.Info, "server actually store length: %v\n", cl)
	//
	//// 修正文件指针
	_, e = f.Seek(cl, io.SeekStart)
	if e != nil {
		return e
	}
	r2, e := http.NewRequest(http.MethodPut, api+location1, nil)
	r2.Body = f
	r2.Header.Set("Range", fmt.Sprintf("bytes=%v-", cl))
	utils.Log.Printf(utils.Info, "do second put...\n")
	rs2, e := client.Do(r2)
	if e != nil {
		return e
	}
	if rs2.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", rs2.Status)
	}
	utils.Log.Printf(utils.Info, "second put over\n")
	return nil
}

func getBigFile(file string, size int64) error {
	r, _ := http.NewRequest(http.MethodGet, apihost+file, nil)
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-", size))
	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("size get %s", response.Status)
	}
	ran := response.Header.Get("content-range")
	fmt.Println(ran)
	data_1, _ := os.Create("data_1")
	defer data_1.Close()
	io.Copy(data_1, response.Body)

	// 第二次 获取全部数据
	r.Header.Set("Range", "bytes=0-")
	response, err = client.Do(r)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("all get %s", response.Status)
	}
	ran = response.Header.Get("content-range")
	fmt.Println(ran)		// 应该是空的
	data_all, _ := os.Create("data_all")
	defer data_all.Close()
	io.Copy(data_all, response.Body)
	return nil
}