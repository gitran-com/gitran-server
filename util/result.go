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
	Code:    400,
	Msg:     "bad request",
}

//Result404 is 404 JSON response
var Result404 = Result{
	Success: false,
	Code:    404,
	Msg:     "not found",
}

//ResultInvalidToken is JSON response when token is invalid
var ResultInvalidToken = Result{
	Success: false,
	Code:    401,
	Msg:     "invalid token",
}
