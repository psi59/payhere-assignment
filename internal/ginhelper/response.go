package ginhelper

type Response struct {
	Meta ResponseMeta `json:"meta"`
	Data any          `json:"data,omitempty"`
}

type ResponseMeta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
