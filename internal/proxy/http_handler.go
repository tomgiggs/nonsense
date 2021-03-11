package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"nonsense/internal/config"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"nonsense/pkg/proto"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PB2JSON(pbMsg proto.Message) (jsonStr string) {
	json_str, err := json.Marshal(pbMsg)
	if err == nil {
		jsonStr = string(json_str)
	}
	return
}
func GetPutJson(req *gin.Context) *simplejson.Json{
	bodyByte := make([]byte,1000)
	_,_ = req.Request.Body.Read(bodyByte)
	js, err := simplejson.NewJson([]byte(bodyByte))

	if err != nil {
		return nil
	}
	return js
}
func GetUserInfo(token string) (*common.JwtTokenInfo){
	if token==""{
		return nil
	}
	info,succ := common.JwtDecry(token)
	if !succ {
		return nil
	}
	return &info
}

// 认证中间件
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo := GetUserInfo(c.Request.Header.Get("Authorization"))
		if userInfo == nil {
			c.String(http.StatusUnauthorized,"token invalid")
			c.Abort()
		}
		c.Set("user_info",userInfo)// 设置变量到Context的key中，可以通过Get()取
	}
}

func PostReq(c *gin.Context,resp proto.Message,err error){
	if err != nil {
		status, _ := status.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code":   status.Code(),
			"err_msg":     status.Message(),
		})
		return
	}
	c.JSON(http.StatusOK,resp)
}

type ApiServerV1 struct {
	conf *config.Access
	rmg *RoomManager

}
func NewApiServerV1(c *config.Access) *ApiServerV1{
	return &ApiServerV1{
		conf: c,
		rmg: NewRoomManager(),
	}
}

func (v1 *ApiServerV1)InitRouter(){
	router := gin.Default()
	router.GET("/login",v1.Login)//websocket收发消息处理
	router.PUT("/user/create",v1.AddUser)
	r1 := router.Group("/v1")
	//r1.Use(AuthMiddleWare())
	{
		r1.GET("/ws",v1.WsClient)//websocket收发消息处理
		r1.GET("/rtc",v1.RtcClient)//websocket收发消息处理
		r1.PUT("/device/register",v1.RegisterDevice)
		r1.DELETE("/device/delete",v1.DeleteDevice)
		r1.GET("/device/list",v1.GetDeviceList)
		//user

		r1.POST("/user/update",v1.UpdateUser)
		r1.GET("/user/get",v1.GetUser)
		r1.GET("/user/group/list",v1.GetUserGroups)

		//group
		r1.GET("/group/create",v1.CreateGroup)
		r1.POST("/group/update",v1.UpdateGroup)
		r1.GET("/group/get",v1.GetGroup)
		r1.DELETE("/group/delete",v1.DeleteGroup)
		//group member
		r1.PUT("/group_member/add",v1.AddGroupMember)
		r1.POST("/group_member/update",v1.UpdateGroupMember)
		r1.DELETE("/group_member/delete",v1.DeleteGroupMember)
		r1.GET("/group_member/list",v1.GetGroupMembers)

	}
	router.Run(v1.conf.HttpAddr)
}

func (v1 *ApiServerV1) Close(){

}
func (v1 *ApiServerV1)Login(c *gin.Context){
	body := GetPutJson(c)
	if body ==nil {
		c.String(http.StatusBadRequest,"nil body")
		return
	}
	appIdStr := body.Get("app_id").MustString()
	appId,_ := strconv.ParseInt(appIdStr, 10, 64)
	userIdStr := body.Get("user_id").MustString()
	userId,_ := strconv.ParseInt(userIdStr,10,64)
	passwd :=  body.Get("password").MustString()
	req := &pb.SignInReq{
		AppId: appId,
		UserId: userId,
		Passwd: passwd,
		ConnId: global.AppConfig.AppId,
		DeviceId: 0,
		UserIp: c.Request.RemoteAddr,
	}
	userBasicInfo,_ := c.Get("user_info")
	resp, err := global.WsDispatch.SignIn(common.GetContext(userBasicInfo.(common.JwtTokenInfo)),req)
	if err != nil {
		status, _ := status.FromError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"err_code":   status.Code(),
			"err_msg":     status.Message(),
		})
		return
	}
	c.JSON(http.StatusOK,resp)
}
func (v1 *ApiServerV1)WsClient(c *gin.Context) {
	if ! c.IsWebsocket() {
		c.String(http.StatusOK, "====not websocket request====")
	}
	w,r := c.Writer,c.Request
	upgrader := websocket.Upgrader{}

	appId, _ := strconv.ParseInt(r.Header.Get(common.CtxAppId), 10, 64)
	userId, _ := strconv.ParseInt(r.Header.Get(common.CtxUserId), 10, 64)
	deviceId, _ := strconv.ParseInt(r.Header.Get(common.CtxDeviceId), 10, 64)
	passwd := r.Header.Get(common.CtxToken)

	if appId == 0 || userId == 0 || deviceId == 0 || passwd == "" {
		s, _ := status.FromError(common.ErrUnauthorized)
		bytes, err := json.Marshal(s.Proto())
		if err != nil {
			common.Sugar.Error(err)
			return
		}
		w.Write(bytes)
		return
	}
	_, err := global.WsDispatch.SignIn(common.ContextWithRequstId(context.TODO()), &pb.SignInReq{
		AppId:    appId,
		UserId:   userId,
		DeviceId: deviceId,
		Passwd:    passwd,
		ConnId: global.AppConfig.SrvDisc.ID,
		UserIp:r.RemoteAddr,
	})

	s, _ := status.FromError(err)
	if s.Code() != codes.OK {
		bytes, err := json.Marshal(s.Proto())
		if err != nil {
			common.Sugar.Error(err)
			return
		}
		w.Write(bytes)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.Sugar.Error(err)
		return
	}

	// 断开这个设备之前的连接
	preCtx := LoadWsClientOnline(userId)
	if preCtx != nil {
		preCtx.DeviceId = 0
	}

	ctx := NewWSConnContext(conn, appId, userId, deviceId)
	StoreWsClientOnline(userId, ctx)
	ctx.DoConn()
}

func (v1 *ApiServerV1)RtcClient(c *gin.Context) {
	if ! c.IsWebsocket() {
		c.String(http.StatusOK, "====not websocket request====")
	}
	w,r := c.Writer,c.Request
	upgrader := websocket.Upgrader{}

	appId, _ := strconv.ParseInt(r.Header.Get(common.CtxAppId), 10, 64)
	userId, _ := strconv.ParseInt(r.Header.Get(common.CtxUserId), 10, 64)
	//passwd := r.Header.Get(common.CtxToken)
	//
	//if appId == 0 || userId == 0 || passwd == "" {
	//	s, _ := status.FromError(common.ErrUnauthorized)
	//	bytes, err := json.Marshal(s.Proto())
	//	if err != nil {
	//		common.Sugar.Error(err)
	//		return
	//	}
	//	w.Write(bytes)
	//	return
	//}
	//_, err := global.WsDispatch.SignIn(common.ContextWithRequstId(context.TODO()), &pb.SignInReq{
	//	AppId:    appId,
	//	UserId:   userId,
	//	Passwd:    passwd,
	//	ConnId: global.AppConfig.SrvDisc.ID,
	//	UserIp:r.RemoteAddr,
	//})
	//
	//s, _ := status.FromError(err)
	//if s.Code() != codes.OK {
	//	bytes, err := json.Marshal(s.Proto())
	//	if err != nil {
	//		common.Sugar.Error(err)
	//		return
	//	}
	//	w.Write(bytes)
	//	return
	//}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.Sugar.Error(err)
		return
	}

	ctx := NewRtcContext(v1.rmg,conn, appId, userId )
	ctx.Serve()
}


func(v1 *ApiServerV1) RegisterDevice(c *gin.Context) {
	body := GetPutJson(c)
	if body ==nil {
		c.String(http.StatusBadRequest,"nil body")
		return
	}
	appIdStr := body.Get("app_id").MustString()
	appId,_ := strconv.ParseInt(appIdStr, 10, 64)
	brand := body.Get("brand").MustString()
	model :=  body.Get("model").MustString()
	sysVer :=  body.Get("system_version").MustString()
	sdkVer := body.Get("sdk_version").MustString()
	req := &pb.RegisterDeviceReq{
		AppId: appId,
		Brand: brand,
		Model: model,
		SystemVersion: sysVer,
		SdkVersion: sdkVer,
		Type: int32(2),
	}
	userBasicInfo,_ := c.Get("user_info")
	remoteSrv := global.GetHttpDispathServer(c.Request.RemoteAddr)
	if remoteSrv ==nil {
		c.String(http.StatusInternalServerError,"远程服务不可用")
		c.Abort()
	}
	resp, err := remoteSrv.RegisterDevice(common.GetContext(userBasicInfo.(common.JwtTokenInfo)),req)
	PostReq(c,resp,err)

}
func(v1 *ApiServerV1) DeleteDevice(c *gin.Context) {
	//deviceId := c.Param("device_id")
	c.String(http.StatusOK, "to be design...")
}
func(v1 *ApiServerV1) GetDeviceList(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")

}

func GetPage(w http.ResponseWriter, r *http.Request) {
	output := make(chan bool, 1)
	errors := hystrix.Go("get_page", func() error {
		_, err := http.Get("https://www.baidu.com/")
		if err == nil {
			output <- true
		}
		return err
	}, func(err2 error) error {//失败返回
		fmt.Println("get page fail")
		return nil
	})

	select {
	case out := <-output:
		log.Printf("success %v", out)// success
	case err := <-errors:
		log.Printf("failed %s", err)// failure
	}
}

func(v1 *ApiServerV1) AddUser(c *gin.Context) {
	body := GetPutJson(c)
	if body ==nil {
		c.String(http.StatusBadRequest,"nil body")
		return
	}
	//appIdStr := body.Get("app_id").MustString()
	//appId,_ := strconv.ParseInt(appIdStr, 10, 64)
	name := body.Get("nick_name").MustString()
	avatar :=  body.Get("avatar_url").MustString()
	extra :=  body.Get("system_version").MustString()
	sexStr := body.Get("sex").MustString()
	gender,_ := strconv.Atoi(sexStr)

	req := &pb.AddUserReq{
		User: &pb.User{
			UserId: 1,
			Nickname: name,
			Sex: int32(gender),
			AvatarUrl: avatar,
			Extra: extra,
		},
	}
	//自定义配置
	//hystrix.ConfigureCommand("chat_api", hystrix.CommandConfig{
	//	Timeout:                500,
	//	MaxConcurrentRequests:  100,
	//	ErrorPercentThreshold:  50,
	//	RequestVolumeThreshold: 3,
	//	SleepWindow:            1000,
	//})
	hystrix.Go("add_user", func() error {
		resp, err := global.WsDispatch.AddUser(common.SimpleContext(),req)
		if err != nil {
			fmt.Println("add user failed")
			return err
		}
		PostReq(c,resp,err)
		return nil
	}, func(e error) error {
		fmt.Println("add user failed:",e)
		return nil
	})


}
func(v1 *ApiServerV1) UpdateUser(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")

}
func(v1 *ApiServerV1) GetUser(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")

}
func(v1 *ApiServerV1) GetUserGroups(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")

}

func(v1 *ApiServerV1) CreateGroup(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}
func(v1 *ApiServerV1) UpdateGroup(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}
func(v1 *ApiServerV1) GetGroup(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}
func(v1 *ApiServerV1) DeleteGroup(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}

func(v1 *ApiServerV1) AddGroupMember(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}
func(v1 *ApiServerV1) UpdateGroupMember(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}
func(v1 *ApiServerV1) DeleteGroupMember(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}
func(v1 *ApiServerV1) GetGroupMembers(c *gin.Context) {
	c.String(http.StatusOK, "to be design...")
}

