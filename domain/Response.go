package domain

type Response struct {
	StatusCode int         `json:"status_code"`
	StatusMsg  string      `json:"status_msg,omitempty"`
	Attach     interface{} `json:"attach"`
}
