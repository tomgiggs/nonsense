package common

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "nonsense/pkg/proto"
	"runtime"
	"strings"
)

const name = "nonsense"

const TypeUrlStack = "type_url_stack"

var (
	ErrUnknown           = status.New(codes.Unknown, "error unknown").Err()                           // 服务器未知错误
	ErrUnauthorized      = newError(pb.ErrCode_EC_UNAUTHORIZED, "error unauthorized")                 // 未登录
	ErrDBError           = newError(pb.ErrCode_EC_DB_ERROR, "error db error")                     // 数据库错误
	ErrCacheError        = newError(pb.ErrCode_EC_CACHE_ERROR, "cache error")                      // 缓存错误
	ErrNotInGroup        = newError(pb.ErrCode_EC_IS_NOT_IN_GROUP, "error not in group")              // 没有在群组内
	ErrDeviceNotBindUser = newError(pb.ErrCode_EC_DEVICE_NOT_BIND_USER, "error device not bind user") // 没有在群组内
	ErrBadRequest        = newError(pb.ErrCode_EC_BAD_REQUEST, "error bad request")                   // 请求参数错误
	ErrUserAlreadyExist  = newError(pb.ErrCode_EC_USER_ALREADY_EXIST, "error user already exist")     // 用户已经存在
	ErrGroupAlreadyExist = newError(pb.ErrCode_EC_GROUP_ALREADY_EXIST, "error group already exist")   // 群组已经存在
	ErrGroupNotExist     = newError(pb.ErrCode_EC_GROUP_NOT_EXIST, "error group not exist")           // 群组不存在
	ErrUserNotExist      = newError(pb.ErrCode_EC_USER_NOT_EXIST, "error user not exist")             // 用户不存在
)


func newError(code pb.ErrCode, message string) error {
	return status.New(codes.Code(code), message).Err()
}
func WrapError(err error) error {
	if err == nil {
		return nil
	}

	s := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: err.Error(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   Str2bytes(stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func WrapRPCError(err error) error {
	if err == nil {
		return nil
	}
	e, _ := status.FromError(err)
	s := &spb.Status{
		Code:    int32(e.Code()),
		Message: e.Message(),
		Details: []*any.Any{
			{
				TypeUrl: TypeUrlStack,
				Value:   Str2bytes(GetErrorStack(e) + " --grpc-- \n" + stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func GetErrorStack(s *status.Status) string {
	pbs := s.Proto()
	for i := range pbs.Details {
		if pbs.Details[i].TypeUrl == TypeUrlStack {
			return Bytes2str(pbs.Details[i].Value)
		}
	}
	return ""
}

// Stack 获取堆栈信息
func stack() string {
	var pc = make([]uintptr, 20)
	n := runtime.Callers(3, pc)

	var build strings.Builder
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		n := strings.Index(file, name)
		if n != -1 {
			s := fmt.Sprintf(" %s:%d \n", file[n:], line)
			build.WriteString(s)
		}
	}
	return build.String()
}
