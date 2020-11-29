package global

import "fmt"

var (
	USER_LOGIN_SERVER_KEY = "nonsense-user-login-server"
	USER_SERVER_MAP_KEY_PREFIX = "nonsense-user-server-map:"
	USER_CHANGE_EVENT_OFFLINE = "offline"
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

