package util

//Result is HTTP JSON result
type Result struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Code    int                    `json:"code"`
	Msg     string                 `json:"message"`
}

//Result400 is 400 JSON response
var Result400 = Result{
	Success: false,
	Data:    nil,
	Code:    400,
	Msg:     "bad request",
}

//Result401 is 401 JSON response
var Result401 = Result{
	Success: false,
	Data: map[string]interface{}{
		"can_refresh": false,
	},
	Code: 401,
	Msg:  "unauthorized",
}

//Result404 is 404 JSON response
var Result404 = Result{
	Success: false,
	Data:    nil,
	Code:    404,
	Msg:     "not found",
}
