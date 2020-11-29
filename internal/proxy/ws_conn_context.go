package proxy

import (
	"context"
	"fmt"
	"io"
	"nonsense/internal/global"
	"nonsense/pkg/grpclib"
	pb "nonsense/pkg/proto"
	"strings"
	"time"
	"nonsense/pkg/common"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

const PreConn = -1 // 设备第二次重连时，标记设备的上一条连接

type WSConnContext struct {
	Conn     *websocket.Conn // websocket连接
	AppId    int64           // AppId
	DeviceId int64           // 设备id
	UserId   int64           // 用户id
}

func NewWSConnContext(conn *websocket.Conn, appId, userId, deviceId int64) *WSConnContext {
	return &WSConnContext{
		Conn:     conn,
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
	}
}


func (c *WSConnContext) DoConn() {
	defer common.RecoverPanic()

	for {
		err := c.Conn.SetReadDeadline(time.Now().Add(12 * time.Minute))
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		c.HandlePackage(data)
	}
}

// 处理请求发包
func (c *WSConnContext) HandlePackage(bytes []byte) {
	var input pb.Input
	err := proto.Unmarshal(bytes, &input)
	if err != nil {
		common.Sugar.Error(err)
		c.Release()
		return
	}

	switch input.Type {
	case pb.PackageType_PT_SYNC:
		c.Sync(input)
	case pb.PackageType_PT_HEARTBEAT:
		c.Heartbeat(input)
	case pb.PackageType_PT_MESSAGE_ACK:
		c.MessageACK(input)
	case pb.PackageType_PT_MESSAGE_SEND:
		c.SendMsg(input)
	default:
		common.Logger.Info("switch other")
	}

}

// 离线消息同步
func (c *WSConnContext) Sync(input pb.Input) {
	var sync pb.SyncInput
	err := proto.Unmarshal(input.Data, &sync)
	if err != nil {
		common.Sugar.Error(err)
		c.Release()
		return
	}

	resp, err := global.WsDispatch.Sync(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &pb.SyncReq{
		AppId:    c.AppId,
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Seq:      sync.Seq,
	})

	var message proto.Message
	if err == nil {
		message = &pb.SyncOutput{Messages: resp.Messages}
	}

	c.Output(pb.PackageType_PT_SYNC, input.RequestId, err, message)
}

// 发送消息
func (c *WSConnContext) SendMsg(input pb.Input) {
	var msg pb.SendMessageReq
	err := proto.Unmarshal(input.Data, &msg)
	if err != nil {
		common.Sugar.Error(err)
		c.Release()
		return
	}

	resp, err := global.WsDispatch.SendMessage(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &msg)
	fmt.Println(resp.ProtoMessage)

	var message proto.Message
	if err == nil {
		message = &pb.SendMessageResp{}
	}

	c.Output(pb.PackageType_PT_SYNC, input.RequestId, err, message)
}

// 心跳
func (c *WSConnContext) Heartbeat(input pb.Input) {
	c.Output(pb.PackageType_PT_HEARTBEAT, input.RequestId, nil, nil)
	common.Sugar.Infow("heartbeat", "device_id", c.DeviceId, "user_id", c.UserId)
}

// 消息回执
func (c *WSConnContext) MessageACK(input pb.Input) {
	var messageACK pb.MessageACKReq
	err := proto.Unmarshal(input.Data, &messageACK)
	if err != nil {
		common.Sugar.Error(err)
		c.Release()
		return
	}

	_, _ = global.WsDispatch.MessageACK(grpclib.ContextWithRequstId(context.TODO(), input.RequestId), &messageACK)
}

func (c *WSConnContext) Output(pt pb.PackageType, requestId int64, err error, message proto.Message) {
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

	err = c.Conn.WriteMessage(websocket.BinaryMessage, outputBytes)
	if err != nil {
		common.Sugar.Error(err)
		return
	}
}

// 处理错误
func (c *WSConnContext) HandleReadErr(err error) {
	common.Logger.Debug("read tcp error：", zap.Int64("app_id", c.AppId), zap.Int64("user_id", c.UserId),
		zap.Int64("device_id", c.DeviceId), zap.Error(err))
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	c.Release()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
}

// 释放TCP连接
func (c *WSConnContext) Release() {
	// 关闭tcp连接
	err := c.Conn.Close()
	if err != nil {
		common.Sugar.Error(err)
	}

	// 通知业务服务器设备下线
	if c.DeviceId != PreConn {
		_, _ = global.WsDispatch.Offline(context.TODO(), &pb.OfflineReq{
			AppId:    c.AppId,
			UserId:   c.UserId,
			DeviceId: c.DeviceId,
		})
	}
}
func LoadWsClientOnline(userId int64) *WSConnContext {
	value, ok := global.WSManager.Load(userId)
	if ok {
		return value.(*WSConnContext)
	}
	return nil
}

func StoreWsClientOnline(userId int64, ctx *WSConnContext) {
	global.WSManager.Store(userId, ctx)
}

func DeleteWsClientOnline(userId int64) {
	global.WSManager.Delete(userId)
}