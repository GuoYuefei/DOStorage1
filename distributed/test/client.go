package test

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
)
const filepath = "./"
var apihost = "http://" + config.ServerInf.LISTEN_ADDRESS + "/objects/"

func Put(file string, ok bool) error {
	if ok {
		return putCorrect(file)
	}
	return putIncorrect(file)
}

func Get(name string) error {
	resp, err := http.Get(apihost + name)
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
