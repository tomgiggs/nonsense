package grpclib

import (
	"context"
	"nonsense/pkg/common"
	"strconv"
	"google.golang.org/grpc/metadata"
)

const (
	CtxAppId     = "app_id"
	CtxUserId    = "user_id"
	CtxDeviceId  = "device_id"
	CtxToken     = "passwd"
	CtxRequestId = "request_id"
)

func ContextWithRequstId(ctx context.Context, requestId int64) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestId, strconv.FormatInt(requestId, 10)))
}

// GetCtxAppId 获取ctx的app_id
func GetCtxRequstId(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}

	requstIds, ok := md[CtxRequestId]
	if !ok && len(requstIds) == 0 {
		return 0
	}
	requstId, err := strconv.ParseInt(requstIds[0], 10, 64)
	if err != nil {
		return 0
	}
	return requstId
}

// GetCtxData 获取ctx的用户数据，依次返回app_id,user_id,device_id
func GetCtxData(ctx context.Context) (int64, int64, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, 0, common.ErrUnauthorized
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
		return 0, 0, 0, common.ErrUnauthorized
	}
	appId, err = strconv.ParseInt(appIdStrs[0], 10, 64)
	if err != nil {
		common.Sugar.Error(err)
		return 0, 0, 0, common.ErrUnauthorized
	}

	userIdStrs, ok := md[CtxUserId]
	if !ok && len(userIdStrs) == 0 {
		return 0, 0, 0, common.ErrUnauthorized
	}
	userId, err = strconv.ParseInt(userIdStrs[0], 10, 64)
	if err != nil {
		common.Sugar.Error(err)
		return 0, 0, 0, common.ErrUnauthorized
	}

	deviceIdStrs, ok := md[CtxDeviceId]
	if !ok && len(deviceIdStrs) == 0 {
		return 0, 0, 0, common.ErrUnauthorized
	}
	deviceId, err = strconv.ParseInt(deviceIdStrs[0], 10, 64)
	if err != nil {
		common.Sugar.Error(err)
		return 0, 0, 0, common.ErrUnauthorized
	}
	return appId, userId, deviceId, nil
}

// GetCtxAppId 获取ctx的app_id
func GetCtxAppId(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, common.ErrUnauthorized
	}

	tokens, ok := md[CtxAppId]
	if !ok && len(tokens) == 0 {
		return 0, common.ErrUnauthorized
	}
	appId, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		common.Sugar.Error(err)
		return 0, common.ErrUnauthorized
	}

	return appId, nil
}

// GetCtxAppId 获取ctx的token
func GetCtxToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", common.ErrUnauthorized
	}

	tokens, ok := md[CtxToken]
	if !ok && len(tokens) == 0 {
		return "", common.ErrUnauthorized
	}

	return tokens[0], nil
}
