package model

//Response is HTTP JSON result
type Response struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Code    int                    `json:"code"`
	Msg     string                 `json:"message"`
}

//Resp400 is 400 JSON response for 'Bad Request'
var Resp400 = Response{
	Success: false,
	Code:    400,
	Msg:     "bad request",
}

//Resp404 is 404 JSON response for 'Not Found'
var Resp404 = Response{
	Success: false,
	Code:    404,
	Msg:     "not found",
}

//RespInvalidToken is JSON response when token is invalid
var RespInvalidToken = Response{
	Success: false,
	Code:    401,
	Msg:     "invalid token",
}
