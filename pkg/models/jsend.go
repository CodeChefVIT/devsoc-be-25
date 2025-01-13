package models

type Response struct {
	Status  string      `json:"status"`
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
