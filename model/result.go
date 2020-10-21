package model

type Result struct {
	Success bool                   `json:"success"`
	Msg     string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
