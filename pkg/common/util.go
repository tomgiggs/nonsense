package common

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

//ip转地址字符串
func InetNtoA(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// 格式化时间
func FormatTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// 将时间串转为Time
func ParseTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

// 时间转毫秒数
func UnixMilliTime(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// 毫秒数转时间
func UnunixMilliTime(unix int64) time.Time {
	return time.Unix(0, unix*1000000)
}
// 恢复panic
func RecoverPanic() {
	err := recover()
	if err != nil {
		Logger.DPanic("panic", zap.Any("panic", err), zap.String("stack", GetStackInfo()))
	}
}

// 获取Panic堆栈信息
func GetStackInfo() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}

func JsonMarshal(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		Logger.Error("json序列化：", zap.Error(err))
	}
	return Bytes2str(bytes)
}
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetLocalIp() (localAddr string){
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		Logger.Error("",zap.Any("err:",err))
		return
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				Logger.Debug("find local ip: ",zap.Any("ip:",ipnet.IP.String()))
				//ipnet.IP.String()
				if !strings.HasSuffix(ipnet.IP.String(),".0.1"){
					localAddr = ipnet.IP.String()
				}
			}
		}
	}
	Logger.Debug("final local ip: ",zap.Any("ip:",localAddr))
	return
}