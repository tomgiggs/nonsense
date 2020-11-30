package proxy

import (
	"context"
	"go.uber.org/zap"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
)

type LogicDispatchServer struct {
	pb.UnimplementedLogicDispatchServer
}

func DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) error {
	userMap := global.UserFdMap[req.UserId]
	if userMap == nil {
		return nil
	}
	for fd, _ := range userMap {
		connection, ok := global.TcpServer.GetConn(fd) // 获取设备对应的TCP连接
		if !ok {
			common.Logger.Warn("GetConn warn", zap.Int64("user_id", req.UserId))
			continue
		}
		global.SendToClient(connection, pb.PackageType_PT_MESSAGE_ACK, common.GetCtxRequstId(ctx), nil, req.Message)
	}
	wsClient := LoadWsClientOnline(req.UserId)
	if wsClient != nil {
		wsClient.Output(pb.PackageType_PT_SYNC, common.GetCtxRequstId(ctx), nil, req.Message)
	}
	return nil
}

// Message 投递消息
func (s *LogicDispatchServer) DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) (*pb.DeliverMessageResp, error) {
	return &pb.DeliverMessageResp{}, DeliverMessage(ctx, req)
}
