package store

import (
	"database/sql"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"nonsense/internal/config"
	"nonsense/pkg/common"
)
var StorageClient *DBClient

type DBClient struct{
	MysqlClient *sql.DB
	RedisClient *redis.Client
	Conf  *config.Access
}

func NewDBClient(c *config.Access)*DBClient {
	client := &DBClient{
		Conf: c,
	}
	client.Init()
	return client
}

func (self *DBClient)Init() {
	DBCli, err := sql.Open("mysql", self.Conf.Storage.MySQL)
	if err != nil {
		panic(err)
	}
	self.MysqlClient=DBCli
	//init redis client
	addr := self.Conf.Storage.Redis
	RedisCli := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
		Password: self.Conf.Storage.RedisPasswd,
	})

	_, err = RedisCli.Ping().Result()
	if err != nil {
		common.Sugar.Error("redis err ")
		panic(err)
	}
	self.RedisClient = RedisCli
}

func (self *DBClient)Close(){
	if self.MysqlClient != nil{
		self.MysqlClient.Close()
	}
	if self.RedisClient != nil{
		self.RedisClient.Close()
	}
}