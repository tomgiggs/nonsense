package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net/http"
	"nonsense/pkg/common"
	"nonsense/pkg/proto"
	"strings"
	"time"
)

func WsClient(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.Sugar.Error(err)
		return
	}

	ctx := NewWSConnContext(conn, 1, 2, 3)
	ctx.DoConn()
}
func RtcClient(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:4000,
		WriteBufferSize:1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.Sugar.Error(err)
		return
	}
	 rmg := NewRoomManager()
	ctx := NewRtcContext(rmg,conn, 1, 2 )
	ctx.Serve()
}
func StartWsServer(){
	http.HandleFunc("/ws", WsClient)
	http.HandleFunc("/rtc",RtcClient)
	if err := http.ListenAndServe("127.0.0.1:16001", nil); err != nil {
		log.Fatal("start http server err:",err)
	}
}


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
	fmt.Printf("raw input is:%v\n",input)
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
	fmt.Printf("send msg: %v\n",msg)
	pbMessages := make([]*pb.MessageItem, 0)
	oneMsg := &pb.MessageItem{
		MessageId:"dddd",
		MessageBody:nil,
		SenderId:2,
	}
	pbMessages = append(pbMessages,oneMsg)
	receivedMsg := &pb.SyncResp{
		Messages:pbMessages,
	}
	c.Output(pb.PackageType_PT_SYNC,"",nil,receivedMsg)
	return

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

}

func (c *WSConnContext) Output(pt pb.PackageType, requestId string, err error, message proto.Message) {
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
}
