package proxy

import (
	"context"
	"fmt"
	"github.com/alberliu/gn"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"nonsense/internal/config"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
	"strconv"
	"time"
)


func InitLogicDispatchClient(addr string) {
	if global.LogicDispatchMap[addr] != nil {
		return
	}
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(DispatcherInterceptor))
	if err != nil {
		common.Sugar.Error(err)
		panic(err)
	}
	global.LogicDispatchMap[addr] = pb.NewLogicDispatchClient(conn)
}
func InitWsDispatchClient(addr string) {
	if global.HttpDispathMap[addr] != nil {
		return
	}
	conn, err := grpc.DialContext(context.TODO(), addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(DispatcherInterceptor))
	if err != nil {
		common.Sugar.Error(err)
		panic(err)
	}
	global.WsDispatch = pb.NewLogicClientExtClient(conn)
	global.HttpDispathMap[addr] = pb.NewLogicClientExtClient(conn)
	global.HttpSrvList = append(global.HttpSrvList, addr)
}

func ServiceDiscover(conf *config.Access) {
	var lastIndex uint64
	config := api.DefaultConfig()
	config.Address = conf.SrvDisc.Addr
	client, err := api.NewClient(config)
	if err != nil {
		common.Sugar.Error("api new client is failed, err:", err)
		return
	}
	services, metainfo, err := client.Health().Service(conf.SrvDisc.ServiceName, conf.SrvDisc.Tag, true, &api.QueryOptions{
		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		logrus.Warn("error retrieving instances from Consul: %v", err)
	}
	lastIndex = metainfo.LastIndex

	for _, service := range services {
		common.Logger.Info("service dis",zap.Any("node.Address:", service.Service.Address), zap.Any("node.Id:",service.Service.ID))
		InitLogicDispatchClient(net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port)))
		InitWsDispatchClient(net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port)))
	}
}

func RegisterService(conf *config.Access)  {
	// 创建连接consul服务配置
	config := api.DefaultConfig()
	localIp := common.GetLocalIp()
	config.Address = conf.SrvDisc.Addr
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)
	registration.ID = conf.SrvDisc.ID//节点id
	registration.Name = conf.SrvDisc.ServiceName//服务名称
	registration.Port = conf.LocalDisPort //节点端口
	registration.Tags = []string{conf.SrvDisc.Tag}
	registration.Address = localIp//节点ip

	// 增加consul健康检查回调函数
	check := new(api.AgentServiceCheck)
	check.TCP = localIp+conf.LocalDispAddr
	check.Timeout = conf.SrvDisc.Timeout
	check.Interval = conf.SrvDisc.CheckInterval
	check.DeregisterCriticalServiceAfter = conf.SrvDisc.DeleteDelay // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
	common.Sugar.Info("服务注册成功")
}


//客户端收消息通道
func StartTCPServer(conf *config.Access) {
	var err error
	global.TcpServer, err = gn.NewServer(conf.TcpPort, &Handler{},
		gn.NewHeaderLenDecoder(2, 254),
		gn.WithTimeout(5*time.Minute, 11*time.Minute),
		gn.WithAcceptGNum(10),
		gn.WithIOGNum(10))
	if err != nil {
		panic(err)
	}
	global.TcpServer.Run()
	common.Sugar.Info("tcp服务已启动")
}

//客户端发消息通道
func StartClientRpcServer(conf *config.Access) {
	defer common.RecoverPanic()
	intListen, err := net.Listen("tcp", conf.ClientRpcAddr)
	if err != nil {
		panic(err)
	}
	intServer := grpc.NewServer(grpc.UnaryInterceptor(ClientReqInterceptor))
	pb.RegisterLogicClientExtServer(intServer, &ClientApiServer{})
	err = intServer.Serve(intListen)
	if err != nil {
		panic(err)
	}
	common.Sugar.Info("rpc接口调用服务已启动")
}

//为其他服务器转发消息
func StartDispatchRPCServer(conf *config.Access) {
	listener, err := net.Listen("tcp", conf.LocalDispAddr)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerInterceptor))
	pb.RegisterLogicDispatchServer(server, &LogicDispatchServer{})
	common.Logger.Debug("rpc转发服务已启动")
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}


func StartHttpServer(conf *config.Access) {
	httpServer := NewApiServerV1(conf)
	httpServer.InitRouter()
	common.Sugar.Info("http 服务已启动")
}
