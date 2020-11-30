package global

import (
	"github.com/alberliu/gn"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/status"
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
	Encoder = gn.NewHeaderLenEncoder(2, 1024)
	WSManager sync.Map
	AppConfig *config.Access
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

