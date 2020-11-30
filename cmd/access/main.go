package main

import (
	"nonsense/internal/config"
	"nonsense/internal/global"
	"nonsense/internal/proxy"
	"nonsense/internal/store"
	"nonsense/pkg/common"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){
	conf,err := config.Init()
	if err !=nil {
		panic(err)
	}
	global.AppConfig = conf
	store.StorageClient= store.NewDBClient(conf)
	store.OpenAdapter(conf)
	store.NewRcache(conf)

	//为其他服务器转发消息到客户端
	go func() {
		defer common.RecoverPanic()
		proxy.StartDispatchRPCServer(conf) //启动服务
	}()

	go func() {
		proxy.RegisterService(conf)        //注册服务
		for{
			proxy.ServiceDiscover(conf)
			time.Sleep(time.Second*30)
		}

	}()
	go func() {
		//rpc客户端消息通道
		proxy.StartClientRpcServer(conf)
	}()
	go func() {
		proxy.StartWSServer(conf)// 客户端ws消息通道
	}()
	// 客户端tcp消息通道
	go func() {
		proxy.StartTCPServer(conf)
	}()

	//进程退出处理
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		common.Sugar.Infof("nonsense get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:

			if 	store.Storage != nil {
				store.Storage.Close()
			}
			if global.TcpServer != nil{
				global.TcpServer.Stop()
				time.Sleep(time.Second)
			}

			return
		case syscall.SIGHUP:

		default:
			return
		}
	}
}
