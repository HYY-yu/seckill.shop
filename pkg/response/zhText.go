package response

var codeText = map[int]string{
	ServerError:        "服务器错误",
	TooManyRequests:    "请求发送过多",
	AuthorizationError: "鉴权失败",
	ParamBindError:     "请检查参数是否在正确",
}
