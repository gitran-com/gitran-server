package model

//Result is HTTP JSON result
type Result struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Code    int                    `json:"code"`
	Msg     string                 `json:"message"`
}

//Result401 is 401 JSON response
var Result401 = Result{
	Success: false,
	Data:    nil,
	Code:    401,
	Msg:     "unauthorized",
}

//Result404 is 404 JSON response
var Result404 = Result{
	Success: false,
	Data:    nil,
	Code:    404,
	Msg:     "not found",
}
