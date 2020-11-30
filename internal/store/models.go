package store

import (
	pb "nonsense/pkg/proto"
	"time"
)

type App struct {
	Id         int64     // AppId
	Name       string    // 名称
	PrivateKey string    // 私钥
	CreateTime time.Time // 创建时间
	UpdateTime time.Time // 更新时间
}

const (
	DeviceOnLine  = 1 // 设备在线
	DeviceOffLine = 0 // 设备离线
)

// Device 设备
type Device struct {
	Id            int64     // 设备id
	DeviceId      int64     // 设备id
	AppId         int64     // app_id
	UserId        int64     // 用户id
	Type          int32     // 设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web
	Brand         string    // 手机厂商
	Model         string    // 机型
	SystemVersion string    // 系统版本
	SDKVersion    string    // SDK版本
	Status        int32     // 在线状态，0：不在线；1：在线
	ConnId        string    // 连接层服务层地址
	UserIp        string    // TCP连接对应的文件描述符
	CreateTime    time.Time // 创建时间
	UpdateTime    time.Time // 更新时间
}

type DeviceToken struct {
	UserId int64
	Token  string
}
type UserSeq struct {
	AppId      int64
	UserId     int64
	ReadSeq    int64
	ReceiveSeq int64
}

// User 账户
type User struct {
	Id       int64  `json:"-"`        // 用户id
	AppId    int64  `json:"app_id"`   // app_id
	UserId   int64  `json:"user_id"`  // 手机号
	Passwd   string `json:"passwd"`   // 密码
	Birthday string `json:"birthday"` //生日
	Mobile   string `json:"birthday"` //手机号
	Email    string `json:"email"`    //邮箱

	Nickname      string    `json:"nickname"`   // 昵称
	Sex           int32     `json:"sex"`        // 性别，1:男；2:女
	AvatarUrl     string    `json:"avatar_url"` // 用户头像
	Extra         string    `json:"extra"`      // 附加属性
	CreateTime    time.Time `json:"-"`          // 创建时间
	LastLoginTime time.Time `json:"-"`          // 更新时间
	WeixinOpenid  string    `json:"weixin_openid"`
	LastLoginIp   string    `json:"last_login_ip"`
	RegisterIp    string    `json:"register_ip"`
}

const (
	GroupTypeGroup    = 1 // 群组
	GroupTypeChatRoom = 2 // 聊天室
)

// Group 群组
type Group struct {
	Id           int64     `json:"-"`            // 群组id
	AppId        int64     `json:"-"`            // appId
	GroupId      int64     `json:"group_id"`     // 群组id
	Name         string    `json:"name"`         // 组名
	Introduction string    `json:"introduction"` // 群简介
	UserNum      int32     `json:"user_num"`     // 群组人数
	Type         int32     `json:"type"`         // 群组类型
	Extra        string    `json:"extra"`        // 附加属性
	CreateTime   time.Time `json:"-"`            // 创建时间
	UpdateTime   time.Time `json:"-"`            // 更新时间
}

type GroupUserUpdate struct {
	GroupId int64  `json:"group_id"` // 群组id
	UserId  int64  `json:"user_id"`  // 用户id
	Label   string `json:"label"`    // 用户标签
	Extra   string `json:"extra"`    // 群组用户附件属性
}

type GroupMember map[int64]GroupUserInfo

// 群组成员信息
type GroupUserInfo struct {
	AppId          int64     `json:"app_id"`           // 应用id
	GroupId        int64     `json:"group_id"`         // 群组id
	UserId         int64     `json:"user_id"`          // 用户id
	Label          string    `json:"label"`            // 用户标签
	GroupUserExtra string    `json:"group_user_extra"` // 群组用户附件属性
	Nickname       string    `json:"name"`             // 昵称
	Sex            int       `json:"sex"`              // 性别,0:位置；1:男；2:女
	AvatarUrl      string    `json:"img"`              // 用户头像
	UserExtra      string    `json:"user_extra"`       // 用户附件属性
	CreateTime     time.Time `json:"create_time"`      // 创建时间
	UpdateTime     time.Time `json:"update_time"`      // 更新时间
}

type Sender struct {
	AppId      int64         // appId
	SenderType pb.SenderType // 发送者类型，1：系统，2：用户，3：业务方
	SenderId   int64         // 发送者id
	DeviceId   int64         // 发送者设备id
}

const (
	MessageObjectTypeUser  = 1 // 用户
	MessageObjectTypeGroup = 2 // 群组
)

// Message 消息
type Message struct {
	Id             int64     // 自增主键
	AppId          int64     // appId
	ObjectType     int       // 所属类型
	ObjectId       int64     // 所属类型id
	MessageId      string    // 请求id
	SenderType     int32     // 发送者类型
	SenderId       int64     // 发送者账户id
	SenderDeviceId int64     // 发送者设备id
	ReceiverType   int32     // 接收者账户id
	ReceiverId     int64     // 接收者id私聊为user_id，群组消息为group_id
	ToUserIds      string    // 需要@的用户id列表，多个用户用，隔开
	Type           int       // 消息类型
	Content        string    // 消息内容
	Seq            int64     // 消息同步序列
	SendTime       time.Time // 消息发送时间
	Status         int32     // 创建时间
}

// 发送消息请求
type SendMessage struct {
	ReceiverType pb.ReceiverType `json:"receiver_type"`
	ReceiverId   int64           `json:"receiver_id"`
	ToUserIds    []int64         `json:"to_user_ids"`
	MessageId    string          `json:"message_id"`
	SendTime     int64           `json:"send_time"`
	MessageBody  struct {
		MessageType    int    `json:"message_type"`
		MessageContent string `json:"-"`
	} `json:"message_body"`
	PbBody *pb.MessageBody `json:"-"`
}
