package gen

const JsonResponseKey = "__gen_JsonResponseKey{}"

type JSON struct {
	Code       int         `json:"code,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	ErrorStack []string    `json:"error_stack,omitempty"`
}

type JsonResponse interface {
	SetCode(int) JsonResponse
	SetMessage(string) JsonResponse
	SetData(interface{}) JsonResponse
	SetErrorStack([]string) JsonResponse
}

func (j *JSON) SetCode(code int) JsonResponse {
	j.Code = code
	return j
}

func (j *JSON) SetMessage(message string) JsonResponse {
	j.Message = message
	return j
}

func (j *JSON) SetData(data interface{}) JsonResponse {
	j.Data = data
	return j
}

func (j *JSON) SetErrorStack(stack []string) JsonResponse {
	j.ErrorStack = stack
	return j
}
