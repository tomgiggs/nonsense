package common

import (
	"fmt"
	"time"
)

var (
	USER_LOGIN_SERVER_KEY = "nonsense-user-login-server"
	USER_SERVER_MAP_KEY_PREFIX = "nonsense-user-server-map:"
	USER_CHANGE_EVENT_OFFLINE = "offline"
	GROUP_CACHE_KEY = "group:"
	GROUP_USER_CACHE_KEY = "group_user:"
	GROUP_CACHE_EXPIRE = time.Minute*5
	USER_CACHE_KEY = "user:"
	USER_CACHE_EXPIRE = time.Minute*5
	APP_CACHE_KEY = "app:"
	APP_CACHE_EXPIRE = time.Hour
	DEVICE_CACHE_KEY = "device:"
	DEVICE_CACHE_EXPIRE = time.Hour
)

var (
	USER_SEQ_KEY_PREFIX ="user_seq:"
	USER_ACK_KEY_PREFIX = "user_ack:"
	GROUP_SEQ_KEY_PREFIX = "group_seq:"
)

var (
	REQ_RESULT_CODE_OK = int32(200)
	REQ_RESULT_CODE_FAIL = int32(400)
	REQ_RESULT_CODE_DB_ERR = int32(401)
	REQ_RESULT_CODE_ILLEGAL= int32(402)
)
var (
	DB_ERROR = fmt.Errorf("query db data failed")
)

const KeyInfo = "ABCDEFGsdwefwwefwfwe323HIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()"