package interfaces

type APIResponse interface {
    SetStatus(status string)
    SetMessage(message string)
    SetData(data interface{})
}

type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func (r *Response) SetStatus(status string) {
    r.Status = status
}

func (r *Response) SetMessage(message string) {
    r.Message = message
}

func (r *Response) SetData(data interface{}) {
    r.Data = data
}

func NewResponse() *Response {
    return &Response{}
}