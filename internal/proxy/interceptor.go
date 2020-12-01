package proxy

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"nonsense/pkg/common"
)

func logPanic(serverName string, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, err *error) {
	p := recover()
	if p != nil {
		common.Logger.Error(serverName+" panic", zap.Any("info", info), zap.Any("ctx", ctx), zap.Any("req", req),
			zap.Any("panic", p), zap.String("stack", common.GetStackInfo()))
		*err = common.ErrUnknown
	}
}


func doClientValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//token, err := common.GetCtxToken(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = service.AuthServiceInst.IsTokenExpire(ctx, token)
	//if err != nil {
	//	return nil, err
	//}

	return handler(ctx, req)
}

// 通用拦截器
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	common.Logger.Debug("server interceptor", zap.Any("requestId", common.GetCtxRequstId(ctx)), zap.Any("req", req), zap.Any("resp", resp))
	return resp, err
}

// 客户端发消息通道拦截器
func ClientReqInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		logPanic("client interceptor", ctx, req, info, &err)
	}()

	common.Logger.Debug("client interceptor", zap.Any("requestId", common.GetCtxRequstId(ctx)),
		zap.Any("method", info.FullMethod), zap.Any("req", req), zap.Error(err))
	//对客户端请求进行token验证
	resp, err = doClientValidate(ctx, req, info, handler)

	s, _ := status.FromError(err)
	if s.Code() != 0 && s.Code() < 1000 {
		md, _ := metadata.FromIncomingContext(ctx)
		common.Logger.Error("client interceptor",zap.Any("request_idd", common.GetCtxRequstId(ctx)), zap.String("method", info.FullMethod),
			zap.Any("md", md), zap.Any("req_params", req), zap.Any("resp", resp), zap.Error(err), zap.String("stack", common.GetErrorStack(s)))
	}
	return
}
func DispatcherInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	return common.WrapRPCError(err)
}
