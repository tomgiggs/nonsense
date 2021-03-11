package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"nonsense/pkg/common"
	"sync"
	"time"
)
const OFFLINE_TIMEOUT = 60
type rtcClientMsg struct {
	Cmd      string `json:"cmd"`
	RoomID   int64 `json:"room_id"`
	UserId int64 `json:"user_id"`
	Msg      string `json:"msg"`
}

type RoomManager struct {
	mlock sync.Mutex
	Rooms map[int64]*Room
	RoomCount int32
}
func NewRoomManager()*RoomManager{
	return &RoomManager{
		Rooms:make(map[int64]*Room),
		mlock:sync.Mutex{},
	}
}

func (rmg *RoomManager)CleanRoom(){
	go func() {
		for{
			time.Sleep(time.Second*10)
			for rid,room := range rmg.Rooms {
				if room.memberCount == 0 {
					rmg.mlock.Lock()
					delete(rmg.Rooms,rid)
					rmg.mlock.Unlock()
				}
			}
		}
	}()
}

func (rmg *RoomManager)GetRoomCount()int32{
	return rmg.RoomCount
}
func (rmg *RoomManager)GetRoom(roomId int64) *Room{
	rmg.mlock.Lock()
	if _,ok := rmg.Rooms[roomId];!ok {
		rmg.Rooms[roomId] = NewRoom(roomId)
	}
	rmg.mlock.Unlock()
	return rmg.Rooms[roomId]
}

type RtcContext struct {
	Conn  *websocket.Conn
	UserId int64
	AppId int64
	LastHeartbeat int64
	RoomMgr *RoomManager
	lock sync.Mutex
}
func NewRtcContext(rmg *RoomManager,conn *websocket.Conn, appId, userId int64)*RtcContext{
	return &RtcContext{
		Conn:conn,
		AppId:appId,
		UserId:userId,
		RoomMgr:rmg,
		lock:sync.Mutex{},

	}
}

func (rtc *RtcContext) Serve(){
	defer common.RecoverPanic()

	for {
		//err := rtc.Conn.SetReadDeadline(time.Now().Add(time.Minute))
		_, data, err := rtc.Conn.ReadMessage()
		if err != nil {
			common.Sugar.Error("read rtc ws msg error:",err,data)
			return
		}

		rtc.HandlePackage(data)
	}
}

func (rtc *RtcContext)HandlePackage(bytes []byte){
	msg := &rtcClientMsg{}
	err := json.Unmarshal(bytes,msg)
	if err != nil {
		fmt.Println("json unmarshal error:",err)
		return
	}
	fmt.Println("client msg is:",msg)
	switch msg.Cmd {
	case "enterRoom":
		rtc.EnterRoom(msg)
	case "leaveRoom":
		rtc.LeaveRoom(msg)
	case "send":
		rtc.SendMsg(msg)
	case "heartbeat":
		rtc.LastHeartbeat = time.Now().Unix()
	default:
		common.Sugar.Info("unknown rtc cmd type:",msg.Cmd)
	}
}

func (rtc *RtcContext)EnterRoom(msg *rtcClientMsg){
	room :=	rtc.RoomMgr.GetRoom(msg.RoomID)
	room.AddMember(msg.UserId,rtc)
}

func (rtc *RtcContext)LeaveRoom(msg *rtcClientMsg){
	room := rtc.RoomMgr.GetRoom(msg.RoomID)
	room.RemoveMember(msg.UserId)
}

func (rtc *RtcContext)SendMsg(msg *rtcClientMsg){//为什么一发送完消息连接就断开了呢？
	rtc.lock.Lock()

	outputBytes,err := json.Marshal(msg)
	if err != nil {
		common.Sugar.Error("json marshal error:",err)
	}
	if rtc.Conn == nil {
		common.Sugar.Error("rtc  client nil")
		return
	}
	//err := rtc.Conn.WriteJSON(msg)
	err = rtc.Conn.WriteMessage(websocket.TextMessage, outputBytes)
	if err != nil {
		common.Sugar.Error("writejson to client error:",err)
	}
	rtc.lock.Unlock()
}

func (rtc *RtcContext)Close() (err error){
	if rtc.Conn != nil {
		err = rtc.Conn.Close()
	}

	if err != nil {
		common.Sugar.Error("close rtc conn error:",err)
	}
	return
}

//room
type Room struct {
	Manager *RoomManager
	roomId int64
	members map[int64]*RtcContext
	memberCount int32
	lock    sync.Mutex
	maxMember int32
	regTime	time.Time
}

func NewRoom(roomId int64)*Room	{
	room:= &Room{
		members:make(map[int64]*RtcContext),
		maxMember:30,
		roomId:roomId,
		lock:sync.Mutex{},
	}
	//go func() {
	//	for {
	//		time.Sleep(OFFLINE_TIMEOUT)
	//		room.CleanOffline()
	//	}
	//}()
	return room
}

func(room *Room)AddMember(userId int64,rtcClient *RtcContext){
	if room.memberCount >= room.maxMember {
		return
	}
	_,ok := room.members[userId]
	if ok {
		room.RemoveMember(userId)
		return
	}
	room.lock.Lock()
	room.members[userId]=rtcClient
	room.memberCount += 1
	room.lock.Unlock()
	enterMsg := &rtcClientMsg{
		Cmd:"enter",
		UserId:userId,
		RoomID:room.roomId,
		Msg:"",
	}
	room.Broadcast(enterMsg)
}

func(room *Room)RemoveMember(userId int64){
	_,ok := room.members[userId]
	if !ok {
		return
	}
	room.lock.Lock()
	delete(room.members,userId)
	room.memberCount -= 1
	room.lock.Unlock()
	leaveMsg := &rtcClientMsg{
		Cmd:"leave",
		UserId:userId,
		RoomID:room.roomId,
		Msg:"",
	}
	room.Broadcast(leaveMsg)
}

func(room *Room)Broadcast(msg *rtcClientMsg){
	for _,client := range room.members {
		if client == nil {
			return
		}
		client.SendMsg(msg)
	}
}

func(room *Room)CleanOffline(){
	curUnix := time.Now().Unix()
	for uid,client := range room.members {
		if curUnix - client.LastHeartbeat > OFFLINE_TIMEOUT{
			room.RemoveMember(uid)
			client.Close()
		}
	}
}

func(room *Room)EmptyRoom(){
	leaveMsg := &rtcClientMsg{
		Cmd:"empty",
		RoomID:room.roomId,
		Msg:"",
	}
	room.Broadcast(leaveMsg)
}
