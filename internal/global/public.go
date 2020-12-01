package global

import (
	"github.com/alberliu/gn"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"
	"math/big"
	"net"
	"nonsense/internal/config"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
	"sync"
)

var (
	TcpServer *gn.Server
	UserFdMap = make(map[int64]map[int32]int32,0)
	LogicDispatchMap = make(map[string]pb.LogicDispatchClient,0)
	WsDispatch pb.LogicClientExtClient
	HttpDispathMap = make(map[string]pb.LogicClientExtClient,0)
	HttpSrvList = make([]string,0)
	Encoder = gn.NewHeaderLenEncoder(2, 1024)
	WSManager sync.Map
	AppConfig *config.Access
	UserTokenInfos =make(map[int64]*common.JwtTokenInfo)
)

func SendToClient(c *gn.Conn, pt pb.PackageType, requestId string, err error, message proto.Message) {
	var output = pb.Output{
		Type:      pt,
		RequestId: requestId,
	}

	if err != nil {
		status, _ := status.FromError(err)
		output.Code = int32(status.Code())
		output.Message = status.Message()
	}

	if message != nil {
		msgBytes, err := proto.Marshal(message)
		if err != nil {
			common.Sugar.Error(err)
			return
		}
		output.Data = msgBytes
	}

	outputBytes, err := proto.Marshal(&output)
	if err != nil {
		common.Sugar.Error(err)
		return
	}
	err = Encoder.EncodeToFD(c.GetFd(), outputBytes)
	if err != nil {
		common.Sugar.Error(err)
		return
	}
}


func InetAtoN(ip string) int64 {
	validIp := net.ParseIP(ip).To4()
	if validIp==nil{
		return 0
	}
	ret := big.NewInt(0)
	ret.SetBytes(validIp)
	return ret.Int64()
}

//负载均衡,按ip路由,暂时随机选吧
func GetHttpDispathServer(userIp string) pb.LogicClientExtClient{
	if len(HttpDispathMap)==0{
		return nil
	}
	idx := InetAtoN(userIp)%int64(len(HttpDispathMap))
	if len(HttpSrvList)<int(idx)	{
		return nil
	}
	return HttpDispathMap[HttpSrvList[idx]]
	//i := int64(0)
	//for _,v := range HttpDispathMap {
	//	if i ==idx{
	//		return v
	//	}
	//	i ++
	//}
	//return nil
}
