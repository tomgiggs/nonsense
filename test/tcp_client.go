package main

import (
	"fmt"
	util2 "github.com/alberliu/gn/test/util"
	"github.com/golang/protobuf/proto"
	"net"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
	"strconv"
	"time"
)



type TcpClient struct {
	AppId    int64
	UserId   int64
	DeviceId int64
	Seq      int64
	Passwd   string
	codec    *util2.Codec
}

func StartTcpClient() {
	client := TcpClient{
		AppId: 1,
		UserId: 2,
		DeviceId: 3,
		Seq: 0,
		Passwd: "ssss",
	}
	client.Start()
	select {}
}
func (c *TcpClient) Output(pt pb.PackageType, requestId string, message proto.Message) {
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

	inputByf, err := proto.Marshal(&input)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = c.codec.Conn.Write(util2.Encode(inputByf))
	if err != nil {
		fmt.Println(err)
	}
}

func (c *TcpClient) Start() {
	connect, err := net.Dial("tcp", "localhost:18001")
	if err != nil {
		fmt.Println(err)
		return
	}

	c.codec = util2.NewCodec(connect)

	c.SignIn()
	c.SendMsg()
	c.SyncTrigger()
	go c.Heartbeat()
	go c.Receive()
}

func (c *TcpClient) SignIn() {
	signIn := pb.SignInInput{
		AppId:    c.AppId,
		UserId:   c.UserId,
		DeviceId: c.DeviceId,
		Passwd:    c.Passwd,
	}
	c.Output(pb.PackageType_PT_SIGN_IN, "", &signIn)
}

func (c *TcpClient) SendMsg() {
	body:= &pb.SendMessageReq{
		//MessageId: "11116",
		ReceiverId: 3,
		ToUserIds: []int64{3},
		ReceiverType: pb.ReceiverType_RT_USER,
		IsPersist: true,
		SendTime: time.Now().Unix(),
		MessageBody: &pb.MessageBody{
			MessageType: pb.MessageType_MT_TEXT,
			MessageContent: &pb.MessageContent{
				Content: &pb.MessageContent_Text{
					Text: &pb.Text{
						Text: "hello this is test msg"+strconv.FormatInt(time.Now().Unix(),10),
					},
				},
			},
		},
	}
	bytes, err := proto.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := pb.Input{
		Type: pb.PackageType_PT_MESSAGE_SEND,
		RequestId: "26",
		Data: bytes,
	}
	inputByf, err := proto.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = c.codec.Conn.Write(util2.Encode(inputByf))
	if err != nil {
		fmt.Println(err)
	}
}

func (c *TcpClient) SyncTrigger() {
	c.Output(pb.PackageType_PT_SYNC, "", &pb.SyncInput{Seq: c.Seq})
}

func (c *TcpClient) Heartbeat() {
	ticker := time.NewTicker(time.Second * 30)
	for range ticker.C {
		c.Output(pb.PackageType_PT_HEARTBEAT, "", nil)
	}
}

func (c *TcpClient) Receive() {
	for {
		_, err := c.codec.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			bytes, ok, err := c.codec.Decode()
			if err != nil {
				fmt.Println(err)
				return
			}

			if ok {
				c.HandlePackage(bytes)
				continue
			}
			break
		}
		time.Sleep(time.Second*3)
	}
}

func (c *TcpClient) HandlePackage(bytes []byte) {
	var output pb.Output
	err := proto.Unmarshal(bytes, &output)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch output.Type {
	case pb.PackageType_PT_SIGN_IN:
		var signInRsp pb.SignInResp
		proto.Unmarshal(output.Data,&signInRsp)
		fmt.Printf(" signin info :%v",signInRsp)
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
		for _, msg := range syncResp.Messages {
			fmt.Printf("消息：发送者类型：%d 发送者id：%d 请求id：%s 接收者类型：%d 接收者id：%d  消息内容：%+v seq：%d \n",
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
		c.Output(pb.PackageType_PT_MESSAGE_ACK, output.RequestId, &ack)
		fmt.Println("离线消息同步结束------")
	default:
		fmt.Println("switch other")
	}
}
