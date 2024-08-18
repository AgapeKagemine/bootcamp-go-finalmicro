package handler

type Response struct {
	StstusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Payload    any    `json:"payload"`
}
