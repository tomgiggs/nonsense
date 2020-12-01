package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strconv"
)

type StorageConf struct {
	MySQL                  string
	Redis                string
	RedisPasswd			   string
	Name   string
}

type LogConf struct {
	Destination string
	Level string
}

type ConsulConf struct {
	Addr  string
	ServiceName string
	Tag  string
	Timeout string
	CheckInterval string
	DeleteDelay  string
	ID	string
}
type Access struct {
	AppId         string
	HttpAddr      string
	TcpPort       int
	TcpAddr       string
	LocalDispAddr string
	LocalDisPort  int
	ClientRpcAddr string
	SrvDisc       *ConsulConf
	Storage       *StorageConf
	LogConfig     *LogConf
}
var (
	confPath  string
	region    string
	env 		string
	host      string
	debugLevel     int64
	AccessConf *Access
)
func init() {
	var (
		defaultHost, _    = os.Hostname()
		defDebug, _   = strconv.ParseInt(os.Getenv("LogLevel"),10,64)
	)
	flag.StringVar(&confPath, "conf", "access-example.toml", "default config path.")
	flag.StringVar(&region, "region", os.Getenv("region"), "avaliable region. or use REGION env variable, value: sh etc.")
	flag.StringVar(&env, "access.env", os.Getenv("access_env"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	flag.StringVar(&host, "host", defaultHost, "machine hostname. or use default machine hostname.")
	flag.Int64Var(&debugLevel, "debug", defDebug, "server debug. or use DEBUG env variable, value: true/false etc.")
	flag.Parse()
}

func Init()(accessConf *Access,err error){
	fmt.Println(confPath)
	accessConf = Default()
	_, err = toml.DecodeFile(confPath, &accessConf)
	return
}
func Default() *Access {
	return &Access{
		AppId:         "",
		HttpAddr:      ":18080",
		TcpAddr:       ":18081",
		TcpPort:       18000,
		ClientRpcAddr: "",
		LocalDisPort:  18002,
		LocalDispAddr: "18082",
		Storage: &StorageConf{
			MySQL: "root:123@tcp(localhost:3306)/gim?charset=utf8&parseTime=true",
			Redis: "127.0.0.1:6379",
			RedisPasswd: "",
		},
		SrvDisc: &ConsulConf{
			Addr: "127.0.0.1:8500",
			ServiceName: "nonsense-conn-server",
			ID: "nonsense-conn-",
			Tag: "",
			Timeout: "5s",
			CheckInterval: "5s",
			DeleteDelay: "30s",
		},
		LogConfig: &LogConf{
			Level: "debug",
			Destination: "console",
		},
	}
}
