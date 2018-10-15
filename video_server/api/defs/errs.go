package defs

type Err struct {
	Error    string `json:"error"`
	ErroCode string `json:"error_code"` //error_code系统原生,|status_code是http的,2个是不同的
}

type ErrResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErrResponse{HttpSC: 400, Error: Err{Error: "Request body is not correct", ErroCode: "001"}}
	ErrorNotAuthUser            = ErrResponse{HttpSC: 401, Error: Err{Error: "User authentication failed", ErroCode: "002"}}
	ErrorDBError                = ErrResponse{HttpSC: 500, Error: Err{Error: "DB ops failed", ErroCode: "003"}}
	ErrorInternalFaults         = ErrResponse{HttpSC: 500, Error: Err{Error: "Internal service error", ErroCode: "004"}}
)
