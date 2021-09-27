package model

//Response is HTTP JSON result
type Response struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Code    int                    `json:"code"`
	Msg     string                 `json:"message"`
}

//Resp400 is for 'Bad Request'
var Resp400 = Response{
	Success: false,
	Code:    400,
	Msg:     "bad request",
}

//Resp403 is for 'Forbidden'
var Resp403 = Response{
	Success: false,
	Code:    403,
	Msg:     "forbidden",
}

//Resp404 is for 'Not Found'
var Resp404 = Response{
	Success: false,
	Code:    404,
	Msg:     "not found",
}

//RespInvalidToken when token is invalid
var RespInvalidToken = Response{
	Success: false,
	Code:    401,
	Msg:     "invalid token",
}
