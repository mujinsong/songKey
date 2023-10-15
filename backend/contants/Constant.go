package contants

const (
	RELATION_MATCHER_NAME = "*domain.RelationMatcher"
	RELATION_QUERY_NAME   = "*domain.RelationQuery"

	SUCCESS                           = 200
	ERROR                             = 500
	INVALID_PARAMS                    = 400
	ERROR_AUTH_INSUFFICIENT_AUTHORITY = 20005
	ERROR_READ_FILE                   = 20006
	ERROR_SEND_EMAIL                  = 20007
	ERROR_CALL_API                    = 20008
	ERROR_UNMARSHAL_JSON              = 20009
	ERROR_DATABASE                    = 30001
	ERROR_OSS                         = 40001
)

var MsgFlags = map[int]string{
	SUCCESS:        "success",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_AUTH_INSUFFICIENT_AUTHORITY: "权限不足",
	ERROR_READ_FILE:                   "读文件失败",
	ERROR_SEND_EMAIL:                  "发送邮件失败",
	ERROR_CALL_API:                    "调用接口失败",
	ERROR_UNMARSHAL_JSON:              "解码JSON失败",

	ERROR_DATABASE: "数据库操作出错，请重试",

	ERROR_OSS: "OSS配置错误",
}
