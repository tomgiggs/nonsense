package proxy

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"nonsense/internal/logic/service"
	"nonsense/pkg/common"
	"nonsense/pkg/grpclib"
)

func logPanic(serverName string, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, err *error) {
	p := recover()
	if p != nil {
		common.Logger.Error(serverName+" panic", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
			zap.Any("panic", p), zap.String("stack", common.GetStackInfo()))
		*err = common.ErrUnknown
	}
}


func doLogicClientExt(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod != "/pb.LogicClientExt/RegisterDevice" {
		appId, userId, _, err := grpclib.GetCtxData(ctx)
		if err != nil {
			return nil, err
		}
		token, err := grpclib.GetCtxToken(ctx)
		if err != nil {
			return nil, err
		}

		err = service.AuthServiceInst.Auth(ctx, appId, userId, "", token)
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

// 通用拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	common.Logger.Debug("interceptor", zap.Any("info", info), zap.Any("req", req), zap.Any("resp", resp))
	return resp, err
}
// 客户端发消息通道拦截器
func ClientReqInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		logPanic("logic_client_ext_interceptor", ctx, req, info, &err)
	}()

	resp, err = doLogicClientExt(ctx, req, info, handler)
	common.Logger.Debug("logic_client_ext_interceptor", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),zap.Any("resp", resp), zap.Error(err))

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		common.Logger.Error("logic_client_ext_interceptor", zap.String("method", info.FullMethod), zap.Any("md", md), zap.Any("req", req),
			zap.Any("resp", resp), zap.Error(err), zap.String("stack", common.GetErrorStack(s)))
	}
	return
}

