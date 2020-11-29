package dao

import (
	jsoniter "github.com/json-iterator/go"
	"nonsense/pkg/common"
	pb "nonsense/pkg/proto"
	"strconv"
	"strings"
)


func FormatUserIds(userId []int64) string {
	build := strings.Builder{}
	for i, v := range userId {
		build.WriteString(strconv.FormatInt(v, 10))
		if i != len(userId)-1 {
			build.WriteString(",")
		}
	}
	return build.String()
}

func UnformatUserIds(userIdStr string) []int64 {
	if userIdStr == "" {
		return []int64{}
	}
	toUserIdStrs := strings.Split(userIdStr, ",")
	toUserIds := make([]int64, 0, len(toUserIdStrs))
	for i := range toUserIdStrs {
		userId, err := strconv.ParseInt(toUserIdStrs[i], 10, 64)
		if err != nil {
			common.Sugar.Error(err)
			continue
		}
		toUserIds = append(toUserIds, userId)
	}
	return toUserIds
}

func MessageToPB(message *Message) *pb.MessageItem {
	return &pb.MessageItem{
		MessageId:      message.MessageId,
		SenderType:     pb.SenderType(message.SenderType),
		SenderId:       message.SenderId,
		SenderDeviceId: message.SenderDeviceId,
		ReceiverType:   pb.ReceiverType(message.ReceiverType),
		ReceiverId:     message.ReceiverId,
		ToUserIds:      UnformatUserIds(message.ToUserIds),
		MessageBody:    NewMessageBody(message.Type, message.Content),
		Seq:            message.Seq,
		SendTime:       message.SendTime.Unix(),
		Status:         pb.MessageStatus(message.Status),
	}
}

func MessagesToPB(messages []Message) []*pb.MessageItem {
	pbMessages := make([]*pb.MessageItem, 0, len(messages))
	for i := range messages {
		pbMessage := MessageToPB(&messages[i])
		if pbMessages != nil {
			pbMessages = append(pbMessages, pbMessage)
		}
	}
	return pbMessages
}

// 将pb协议转化为json字符串
func PBToJsonStr(pbBody *pb.MessageBody) (int, string) {
	if pbBody.MessageType == pb.MessageType_MT_UNKNOWN {
		common.Logger.Error("error message type")
		return 0, ""
	}

	var content interface{}
	switch pbBody.MessageType {
	case pb.MessageType_MT_TEXT:
		content = pbBody.MessageContent.GetText()
	case pb.MessageType_MT_FACE:
		content = pbBody.MessageContent.GetFace()
	case pb.MessageType_MT_VOICE:
		content = pbBody.MessageContent.GetVoice()
	case pb.MessageType_MT_IMAGE:
		content = pbBody.MessageContent.GetImage()
	case pb.MessageType_MT_FILE:
		content = pbBody.MessageContent.GetFile()
	case pb.MessageType_MT_LOCATION:
		content = pbBody.MessageContent.GetLocation()
	case pb.MessageType_MT_COMMAND:
		content = pbBody.MessageContent.GetCommand()
	case pb.MessageType_MT_CUSTOM:
		content = pbBody.MessageContent.GetCustom()
	}

	bytes, err := jsoniter.Marshal(content)
	if err != nil {
		common.Sugar.Error(err)
		return 0, ""
	}

	return int(pbBody.MessageType), common.Bytes2str(bytes)
}

// NewMessageBody 创建一个消息体类型
func NewMessageBody(msgType int, msgContent string) *pb.MessageBody {
	content := pb.MessageContent{}
	switch pb.MessageType(msgType) {
	case pb.MessageType_MT_TEXT:
		var text pb.Text
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &text)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Text{Text: &text}
	case pb.MessageType_MT_FACE:
		var face pb.Face
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &face)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Face{Face: &face}
	case pb.MessageType_MT_VOICE:
		var voice pb.Voice
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &voice)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Voice{Voice: &voice}
	case pb.MessageType_MT_IMAGE:
		var image pb.Image
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &image)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Image{Image: &image}
	case pb.MessageType_MT_FILE:
		var file pb.File
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &file)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_File{File: &file}
	case pb.MessageType_MT_LOCATION:
		var location pb.Location
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &location)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Location{Location: &location}
	case pb.MessageType_MT_COMMAND:
		var command pb.Command
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &command)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Command{Command: &command}
	case pb.MessageType_MT_CUSTOM:
		var custom pb.Custom
		err := jsoniter.Unmarshal(common.Str2bytes(msgContent), &custom)
		if err != nil {
			common.Sugar.Error(err)
			return nil
		}
		content.Content = &pb.MessageContent_Custom{Custom: &custom}
	}

	return &pb.MessageBody{
		MessageType:    pb.MessageType(msgType),
		MessageContent: &content,
	}
}

