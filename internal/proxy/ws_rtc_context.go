package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"nonsense/pkg/common"
	"sync"
	"time"
)
const OFFLINE_TIMEOUT = 60//掉线清理时间，单位秒
const(
	RTC_CLIENT_CMD_ENTER_ROOM = "enterRoom"
	RTC_CLIENT_CMD_LEAVE_ROOM = "leaveRoom"
	RTC_CLIENT_CMD_HEART_BEART = "heartBeart"
	RTC_CLIENT_CMD_SEND_MSG = "sendMsg"
	RTC_CLIENT_CMD_CANDIDATE = "candidate"
	RTC_CLIENT_CMD_OFFER = "offer"
	RTC_CLIENT_CMD_ANSWER = "answer"
	RTC_CLIENT_CMD_EMPTY_ROOM = "emptyRoom"
)

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
	}
}

func (rmg *RoomManager)CleanRoom(){
	go func() {
		for{
			time.Sleep(time.Second*10)
			for rid,room := range rmg.Rooms {
				if room.memberCount == 0 {
					rmg.mlock.Lock()
					defer rmg.mlock.Unlock()
					delete(rmg.Rooms,rid)
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
	defer rmg.mlock.Unlock()
	if _,ok := rmg.Rooms[roomId];!ok {
		rmg.Rooms[roomId] = NewRoom(roomId)
	}
	return rmg.Rooms[roomId]
}
func (rmg *RoomManager)DeleteRoom(roomId int64){
	rmg.mlock.Lock()
	defer rmg.mlock.Unlock()
	if _,ok := rmg.Rooms[roomId];!ok {
		return
	}
	delete(rmg.Rooms,roomId)
}
type RtcContext struct {
	Conn  *websocket.Conn
	UserId int64
	AppId int64
	LastHeartbeat int64
	RoomMgr *RoomManager
	lock   sync.Mutex
}
func NewRtcContext(rmg *RoomManager,conn *websocket.Conn, appId, userId int64)*RtcContext{
	return &RtcContext{
		Conn:conn,
		AppId:appId,
		UserId:userId,
		RoomMgr:rmg,
	}
}

func (rtc *RtcContext) Serve(){
	defer common.RecoverPanic()

	for {
		rtc.Conn.SetReadDeadline(time.Now().Add(time.Minute))
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
	case RTC_CLIENT_CMD_ENTER_ROOM:
		rtc.EnterRoom(msg)
	case RTC_CLIENT_CMD_LEAVE_ROOM:
		rtc.LeaveRoom(msg)
	case RTC_CLIENT_CMD_SEND_MSG:
		rtc.SendMsg(msg)
	case RTC_CLIENT_CMD_HEART_BEART:
		rtc.LastHeartbeat = time.Now().Unix()
	case RTC_CLIENT_CMD_CANDIDATE:
		room := rtc.RoomMgr.GetRoom(msg.RoomID)
		room.iceCandidateInfo[msg.UserId] = msg.Msg
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

func (rtc *RtcContext)AddCandidateInfo(msg *rtcClientMsg){
	room := rtc.RoomMgr.GetRoom(msg.RoomID)
	room.AddCandidate(msg)
}

func (rtc *RtcContext)SendMsg(msg *rtcClientMsg){
	rtc.lock.Lock()
	defer rtc.lock.Unlock()
	outputBytes,err := json.Marshal(msg)
	if err != nil {
		common.Sugar.Error("json marshal error:",err)
	}
	if rtc.Conn == nil {
		common.Sugar.Error("rtc  client nil")
		return
	}
	err = rtc.Conn.WriteMessage(websocket.TextMessage, outputBytes)
	if err != nil {
		common.Sugar.Error("writejson to client error:",err)
	}
}

func (rtc *RtcContext)Close(){
	err := rtc.Conn.Close()
	if err != nil {
		common.Sugar.Error("close rtc conn error:",err)
	}
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
	purge bool//自动清理掉线用户
	iceCandidateInfo map[int64]string
}

func NewRoom(roomId int64)*Room	{
	room:= &Room{
		members:make(map[int64]*RtcContext),
		maxMember:30,
		roomId:roomId,
		purge:false,
		iceCandidateInfo:make(map[int64]string),
	}
	room.CleanOffline()
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
	defer room.lock.Unlock()
	room.members[userId]=rtcClient
	room.memberCount += 1
	enterMsg := &rtcClientMsg{
		Cmd:RTC_CLIENT_CMD_ENTER_ROOM,
		UserId:userId,
		RoomID:room.roomId,
	}
	room.Broadcast(enterMsg)
}

func(room *Room)RemoveMember(userId int64){
	_,ok := room.members[userId]
	if !ok {
		return
	}
	room.lock.Lock()
	defer room.lock.Unlock()
	delete(room.members,userId)
	room.memberCount -= 1
	if room.memberCount ==0 {
		room.Manager.DeleteRoom(room.roomId)
	}
	leaveMsg := &rtcClientMsg{
		Cmd:RTC_CLIENT_CMD_LEAVE_ROOM,
		UserId:userId,
		RoomID:room.roomId,
	}
	room.Broadcast(leaveMsg)
}

func(room *Room)AddCandidate(msg *rtcClientMsg){
	room.lock.Lock()
	defer room.lock.Unlock()
	room.iceCandidateInfo[msg.UserId] = msg.Msg
	room.Broadcast(msg)

}

func(room *Room)Broadcast(msg *rtcClientMsg){
	for _,client := range room.members {
		if client == nil {
			return
		}
		client.SendMsg(msg)
	}
}
func(room *Room)EnablePurge(){
	room.purge = true
}

func(room *Room)CleanOffline(){
	go func() {
		for  {
			time.Sleep(OFFLINE_TIMEOUT)
			if !room.purge{
				continue
			}
			curUnix := time.Now().Unix()
			for uid,client := range room.members {
				if curUnix - client.LastHeartbeat > OFFLINE_TIMEOUT{
					room.RemoveMember(uid)
					client.Close()
				}
			}
		}
	}()

}

func(room *Room)EmptyRoom(){
	leaveMsg := &rtcClientMsg{
		Cmd:RTC_CLIENT_CMD_EMPTY_ROOM,
		RoomID:room.roomId,
	}
	room.Broadcast(leaveMsg)
}
