package errors

import "fmt"

var ErrDecodeRequest = New("cannot decode request", 400)
var ErrValidateRequest = New("cannot validate request", 400)
var ErrAssertID = New("cannot assert id", 500)
var ErrWrongToken = New("wrong token", 401)
var ErrNoFileKey = New("no file with key 'file'", 400)

type CoError struct {
	message string
	code    int
}

// Implements error interface
func (ce CoError) Error() string {
	return ce.message
}

func (ce CoError) Code() int {
	return ce.code
}

func New(message string, code int) CoError {
	return CoError{message, code}
}

func Resolve(err error) (string, int) {
	fmt.Println(err.Error())
	switch err.(type) {
	case CoError:
		return err.Error(), err.(CoError).Code()
	default:
		return "Internal error", 500
	}
}
