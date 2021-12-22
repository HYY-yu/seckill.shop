package response

// JsonResponse
// HTTP 服务常用的返回结构
type JsonResponse struct {
	Code    int         `json:"code"`    // 业务码
	Message string      `json:"message"` // 描述信息
	Data    interface{} `json:"data"`
}

// NewResponse 新建一个 JsonResponse
// 此函数保证 JsonResponse.Data 不为 nil
func NewResponse(payload ...interface{}) *JsonResponse {
	responseData := interface{}(nil)
	if len(payload) > 0 && payload[0] != nil {
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

// 系统级错误码定义
const (
	ServerError        = 10001
	TooManyRequests    = 10002
	AuthorizationError = 10003
	ParamBindError     = 10004
)

// Text 注册表转换
func Text(code int) string {
	return codeText[code]
}
