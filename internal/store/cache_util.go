package store

import (
	"encoding/json"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"strconv"
)

//user discovery
func RefreshOnlineUserServer(){
	SubOnlineUserChange(global.USER_LOGIN_SERVER_KEY, func(info string) {
		changeInfo := new(UserChangeInfo)
		json.Unmarshal([]byte(info),&changeInfo)
		if changeInfo.Event == global.USER_CHANGE_EVENT_OFFLINE{
			if _,ok := UserServerMap[changeInfo.Uid][changeInfo.SrvId];ok{
				delete(UserServerMap[changeInfo.Uid],changeInfo.SrvId)
			}
			return
		}

		UserServerMap[changeInfo.Uid][changeInfo.SrvId] = true
	})
}
func PubOnlineUserChange(data *UserChangeInfo){

	jsonStr,_ := json.Marshal(data)
	err := StorageClient.RedisClient.HSet(global.USER_SERVER_MAP_KEY_PREFIX+strconv.FormatInt(data.Uid,10),data.SrvId,1).Err()
	if err != nil {
		common.Sugar.Error("更新用户所在服务器信息失败:",err)
	}

	sub := StorageClient.RedisClient.Subscribe(global.USER_LOGIN_SERVER_KEY)
	defer sub.Close()
	err = StorageClient.RedisClient.Publish("gim-user-login-server",string(jsonStr)).Err()
	if err != nil {
		common.Sugar.Error("发布用户更新消息失败:",err)
	}
}
func SubOnlineUserChange(channel string, callback func(info string)){

	sub := StorageClient.RedisClient.Subscribe(channel)
	defer sub.Close()
	for msg := range sub.Channel(){
		callback(msg.Payload)
	}
}

type UserServerMapCache map[int64]map[string]bool

var UserServerMap = make(UserServerMapCache,0)

type UserChangeInfo struct {
	Event string
	Uid  int64
	SrvId string
}
