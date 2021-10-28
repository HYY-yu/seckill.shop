package response

type JsonResponse struct {
	Code    int         `json:"code"`    // 业务码
	Message string      `json:"message"` // 描述信息
	Data    interface{} `json:"data"`
}

func NewResponse(payload ...interface{}) *JsonResponse {
	responseData := interface{}(nil)
	if len(payload) > 0 {
		responseData = payload[0]
	} else {
		responseData = make(map[string]interface{})
	}

	return &JsonResponse{
		Code:    0,
		Message: "",
		Data:    responseData,
	}
}

const (
	// 系统级错误码
	ServerError        = 10001
	TooManyRequests    = 10002
	AuthorizationError = 10003
	ParamBindError     = 10004
)

func Text(code int) string {
	return codeText[code]
}
