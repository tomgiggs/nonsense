package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

func StartWsclient() {
	client := TcpClient{
		AppId: 1,
		UserId: 2,
		DeviceId: 3,
		Seq: 1,
	}
	client.Start()
	select {}
}


type WSClient struct {
	AppId    int64
	UserId   int64
	DeviceId int64
	Seq      int64
	Passwd   string
	Conn     *websocket.Conn
}

func (c *WSClient) Start() {
	u := url.URL{Scheme: "ws", Host: "localhost:8081", Path: "/ws"}

	header := http.Header{}
	header.Set(common.CtxAppId, strconv.FormatInt(c.AppId, 10))
	header.Set(common.CtxUserId, strconv.FormatInt(c.UserId, 10))
	header.Set(common.CtxDeviceId, strconv.FormatInt(c.DeviceId, 10))

	header.Set(common.CtxToken, c.Passwd)

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		fmt.Println(err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))
	c.Conn = conn

	c.SyncTrigger()
	go c.Heartbeat()
	go c.Receive()
}

func (c *WSClient) Output(pt pb.PackageType, requestId string, message proto.Message) {
	var input = pb.Input{
		Type:      pt,
		RequestId: requestId,
	}
	if message != nil {
		bytes, err := proto.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}

		input.Data = bytes
	}

	writeBytes, err := proto.Marshal(&input)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.Conn.WriteMessage(websocket.BinaryMessage, writeBytes)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *WSClient) SyncTrigger() {
	c.Output(pb.PackageType_PT_SYNC, "", &pb.SyncInput{Seq: c.Seq})
}

func (c *WSClient) Heartbeat() {
	ticker := time.NewTicker(time.Minute * 4)
	for range ticker.C {
		c.Output(pb.PackageType_PT_HEARTBEAT, "", nil)
		fmt.Println("心跳发送")
	}
}

func (c *WSClient) Receive() {
	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		c.HandlePackage(bytes)
	}
}

func (c *WSClient) HandlePackage(bytes []byte) {
	var output pb.Output
	err := proto.Unmarshal(bytes, &output)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch output.Type {
	case pb.PackageType_PT_HEARTBEAT:
		fmt.Println("心跳响应")
	case pb.PackageType_PT_SYNC:
		fmt.Println("离线消息同步开始------")
		syncResp := pb.SyncOutput{}
		err := proto.Unmarshal(output.Data, &syncResp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("离线消息同步响应:code", output.Code, "message:", output.Message)
		fmt.Printf("%+v \n", output)
		for _, msg := range syncResp.Messages {
			fmt.Printf("消息：发送者类型：%d 发送者id：%d 请求id：%d 接收者类型：%d 接收者id：%d  消息内容：%+v seq：%d \n",
				msg.SenderType, msg.SenderId, msg.MessageId, msg.ReceiverType, msg.ReceiverId, msg.MessageBody.MessageContent, msg.Seq)
			c.Seq = msg.Seq
		}

		ack := pb.MessageACKReq{
			AppId: c.AppId,
			GroupId: 0,
			UserId: c.UserId,
			Seq:   c.Seq,
			ReceiveTime: common.UnixMilliTime(time.Now()),
		}
		c.Output(pb.PackageType_PT_SYNC, output.RequestId, &ack)
		fmt.Println("离线消息同步结束------")
	default:
		fmt.Println("switch other")
	}
}
