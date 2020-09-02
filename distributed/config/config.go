package config

import (
	"flag"
	"github.com/GuoYuefei/DOStorage1/distributed/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

/**
	如果是相对位置，则是相对执行文件的位置
 */
var ConfigFile string

var Pub *SPub
var ServerInf *SInf
var ServerData *SData

const usage string = "-c configFilePath"

type ServerType = int
const (
	TypeSInf = iota
	TypeSData
)

var ObjectRoot string
var TempRoot string

func Init() {
	ObjectRoot = filepath.Join(ServerData.STORAGE_ROOT, "objects")
	TempRoot = filepath.Join(ServerData.STORAGE_ROOT, "temp")
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

}

// serverType
// TypeSInf interface server
// TypeSData data server
func ConfigParse(serverType ServerType) {

	switch serverType {
	case TypeSInf:
		interfaceConfigFile, err := ioutil.ReadFile(path.Join(ConfigFile))
		if err == nil {
			yaml.Unmarshal(interfaceConfigFile, ServerInf)
		}
		utils.Log.Printf(utils.Debug, "after read config file, %v\n", ServerInf, Pub)
	case TypeSData:
		dataConfigFile, err := ioutil.ReadFile(path.Join(ConfigFile))
		if err == nil {
			yaml.Unmarshal(dataConfigFile, ServerData)
		}
		utils.Log.Printf(utils.Debug, "after read config file, %v\n", ServerData, Pub)
	default:
		utils.Log.Println(utils.Warning, "Server Type No MATCH!")
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
	// 最后如果STORAGE_ROOT是相对位置的话，转成绝对路径
	if !filepath.IsAbs(ServerData.STORAGE_ROOT) {
		exePath, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(filepath.Dir(exePath))
		ServerData.STORAGE_ROOT = filepath.Join(path, ServerData.STORAGE_ROOT)
	}

	// 配置结束后初始化
	Init()
}

func Flags(serverType ServerType) {
	switch serverType {
	case TypeSInf:
		flag.StringVar(&ConfigFile, "c", "./config/interface.yml", usage)
		flag.StringVar(&ConfigFile, "-config", "./config/interface.yml", usage)
	case TypeSData:
		flag.StringVar(&ConfigFile, "c", "./config/data.yml", usage)
		flag.StringVar(&ConfigFile, "-config", "./config/data.yml", usage)
	default:
		utils.Log.Println(utils.Warning, "Server Type No MATCH!")
	}
	// 相对位置转换成绝对位置
	if !filepath.IsAbs(ConfigFile) {
		exePath, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(filepath.Dir(exePath))
		ConfigFile = filepath.Join(path, ConfigFile)
	}
	utils.Log.Println(utils.Info, "Use config file: ", ConfigFile)
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