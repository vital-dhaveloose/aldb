package aldberr

import "fmt"

func New(code, msg string, details map[string]interface{}) CanvigaError {
	return CanvigaError{code: code, msg: msg, details: details}
}

func Wrap(err error, code, msg string, details map[string]interface{}) CanvigaError {
	return CanvigaError{code: code, msg: msg, details: details, inner: err}
}

type CanvigaError struct {
	code    string
	msg     string
	details map[string]interface{}
	inner   error
}

func (e CanvigaError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.msg)
}

func (e CanvigaError) Code() string {
	return e.code
}

func (e CanvigaError) Message() string {
	return e.msg
}

func (e CanvigaError) Details() map[string]interface{} {
	return e.details
}

func (e CanvigaError) Det(key string, val interface{}) CanvigaError {
	if e.details == nil {
		e.details = map[string]interface{}{key: val}
	} else {
		e.details[key] = val
	}
	return e
}

type CodeSystem string

const (
	CodeSysHttpStatus = CodeSystem("http-status")
	CodeSystCanviga   = CodeSystem("canviga-error-codes")
)
