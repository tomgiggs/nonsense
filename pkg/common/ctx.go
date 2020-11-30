package common

import (
	"context"
	"github.com/go-basic/uuid"
	"google.golang.org/grpc/metadata"
	"strconv"
)

const (
	CtxAppId     = "app_id"
	CtxUserId    = "user_id"
	CtxDeviceId  = "device_id"
	CtxToken     = "passwd"
	CtxRequestId = "request_id"
)

//设置请求id
func ContextWithRequstId(ctx context.Context) context.Context {
	reqId := uuid.New()
	//md := metadata.Pairs("requestId", reqId)
	//ctx := metadata.NewOutgoingContext(context.Background(), md)
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestId, reqId))
}

// 获取ctx的请求id
func GetCtxRequstId(ctx context.Context) string {
	md,ok := metadata.FromOutgoingContext(ctx)
	//md,ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	requstIds, ok := md[CtxRequestId]
	if !ok && len(requstIds) == 0 {
		return ""
	}
	return requstIds[0]
}

// 获取ctx的用户数据，依次返回app_id,user_id,device_id
func GetCtxData(ctx context.Context) (int64, int64, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, 0, ErrUnauthorized
	}

	var (
		appId    int64
		userId   int64
		deviceId int64
		err      error
	)

	// app_id是必填项
	appIdStrs, ok := md[CtxAppId]
	if !ok && len(appIdStrs) == 0 {
		return 0, 0, 0, ErrUnauthorized
	}
	appId, err = strconv.ParseInt(appIdStrs[0], 10, 64)
	if err != nil {
		Sugar.Error(err)
		return 0, 0, 0, ErrUnauthorized
	}

	userIdStrs, ok := md[CtxUserId]
	if !ok && len(userIdStrs) == 0 {
		return 0, 0, 0, ErrUnauthorized
	}
	userId, err = strconv.ParseInt(userIdStrs[0], 10, 64)
	if err != nil {
		Sugar.Error(err)
		return 0, 0, 0, ErrUnauthorized
	}

	deviceIdStrs, ok := md[CtxDeviceId]
	if !ok && len(deviceIdStrs) == 0 {
		return 0, 0, 0, ErrUnauthorized
	}
	deviceId, err = strconv.ParseInt(deviceIdStrs[0], 10, 64)
	if err != nil {
		Sugar.Error(err)
		return 0, 0, 0, ErrUnauthorized
	}
	return appId, userId, deviceId, nil
}

// 获取ctx的app_id
func GetCtxAppId(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, ErrUnauthorized
	}

	tokens, ok := md[CtxAppId]
	if !ok && len(tokens) == 0 {
		return 0, ErrUnauthorized
	}
	appId, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		Sugar.Error(err)
		return 0, ErrUnauthorized
	}

	return appId, nil
}

// 获取ctx的token
func GetCtxToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrUnauthorized
	}

	tokens, ok := md[CtxToken]
	if !ok && len(tokens) == 0 {
		return "", ErrUnauthorized
	}

	return tokens[0], nil
}
