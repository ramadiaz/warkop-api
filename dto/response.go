package dto

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Error   string      `json:"error,omitempty"`
}
