package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var Pub *SPub
var ServerInf *SInf
var ServerData *SData

func init() {

	Pub = &SPub{
		ES_SERVER: "localhost:9200",
		RABBITMQ_SERVER: "amqp://test:test@localhost:5672/",
	}
	ServerInf = &SInf{
		LISTEN_ADDRESS: "localhost:23333",
	}
	ServerInf.SPub = Pub
	ServerData = &SData{
		LISTEN_ADDRESS: "localhost:23334",
		STORAGE_ROOT: "./data",
	}
	ServerData.SPub = Pub

	// 配置文件优先级在环境变量优先级之下
	dataconfigfile, err := ioutil.ReadFile("./data.yml")
	if err == nil {
		yaml.Unmarshal(dataconfigfile, ServerData)
	}
	interfaceConfigFile, err := ioutil.ReadFile("./interface.yml")
	if err == nil {
		yaml.Unmarshal(interfaceConfigFile, ServerInf)
	}

	if os.Getenv("RABBITMQ_SERVER") != "" {
		Pub.RABBITMQ_SERVER = os.Getenv("RABBITMQ_SERVER")
	}
	if os.Getenv("ES_SERVER") != "" {
		Pub.ES_SERVER = os.Getenv("ES_SERVER")
	}

	if os.Getenv("LISTEN_ADDRESS") != "" {
		ServerData.LISTEN_ADDRESS = os.Getenv("LISTEN_ADDRESS")
		ServerInf.LISTEN_ADDRESS = ServerData.LISTEN_ADDRESS
	}

	if os.Getenv("STORAGE_ROOT") != "" {
		ServerData.STORAGE_ROOT = os.Getenv("STORAGE_ROOT")
	}

}

type SPub struct {
	ES_SERVER string			`yaml:"es_server"`
	RABBITMQ_SERVER string		`yaml:"rabbitmq_server"`
}

type SInf struct {
	*SPub						`yaml:"public"`
	LISTEN_ADDRESS string		`yaml:"listen_address"`
}

type SData struct {
	*SPub						`yaml:"public"`
	LISTEN_ADDRESS string		`yaml:"listen_address"`
	STORAGE_ROOT string			`yaml:"storage_root"`
}