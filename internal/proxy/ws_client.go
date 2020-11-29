package proxy

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"nonsense/internal/global"
	"nonsense/pkg/common"
	"nonsense/pkg/grpclib"
	pb "nonsense/pkg/proto"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 65536,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



func WsHandler(w http.ResponseWriter, r *http.Request) {
	appId, _ := strconv.ParseInt(r.Header.Get(grpclib.CtxAppId), 10, 64)
	userId, _ := strconv.ParseInt(r.Header.Get(grpclib.CtxUserId), 10, 64)
	deviceId, _ := strconv.ParseInt(r.Header.Get(grpclib.CtxDeviceId), 10, 64)
	passwd := r.Header.Get(grpclib.CtxToken)
	requestId, _ := strconv.ParseInt(r.Header.Get(grpclib.CtxRequestId), 10, 64)

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
	_, err := global.WsDispatch.SignIn(grpclib.ContextWithRequstId(context.TODO(), requestId), &pb.SignInReq{
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




