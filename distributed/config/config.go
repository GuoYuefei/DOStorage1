package config

import (
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
	ServerData = &SData{
		LISTEN_ADDRESS: "localhost:23334",
		STORAGE_ROOT: "./data",
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

	// todo config


}

type SPub struct {
	ES_SERVER string
	RABBITMQ_SERVER string
}

type SInf struct {
	*SPub
	LISTEN_ADDRESS string
}

type SData struct {
	*SPub
	LISTEN_ADDRESS string
	STORAGE_ROOT string
}