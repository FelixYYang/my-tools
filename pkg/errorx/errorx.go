package errorx

import "strconv"

const (
	UNKNOWN   = -1
	ReqErr    = 1
	Forbidden = 1000
)

// AppError 服务错误类型，指需要返回用户端的错误，包含Code,Msg属性
type AppError interface {
	error
	Code() int16
	Msg() string
	Unwrap() error
}

func New(code int16, msg string) error {
	return &appError{
		code: code,
		msg:  msg,
	}
}

func NewReqErr(msg string) error {
	return &appError{
		code: ReqErr,
		msg:  msg,
	}
}

func Wrap(err error, msg string) error {
	return &appError{
		code:  UNKNOWN,
		msg:   msg,
		cause: err,
	}
}

type appError struct {
	code  int16
	msg   string
	cause error
}

func (a *appError) Error() string {
	return "code " + strconv.Itoa(int(a.code)) + "," + a.msg + ": " + a.cause.Error()
}

func (a *appError) Code() int16 {
	return a.code
}

func (a *appError) Msg() string {
	return a.msg
}

func (a *appError) Unwrap() error {
	return a.cause
}
