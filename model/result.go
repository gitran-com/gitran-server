package model

//Result is HTTP JSON result
type Result struct {
	Success bool                   `json:"success"`
	Msg     string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

//Result401 is 401 JSON response
var Result401 = Result{
	Success: false,
	Msg:     "Unauthorized",
	Data:    nil,
}

//Result404 is 404 JSON response
var Result404 = Result{
	Success: false,
	Msg:     "Not Found",
	Data:    nil,
}
