package gen

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type BusError struct {
	HttpCode int
	ErrCode  int
	Message  string
	Errors   error
}

func NewBusError(httpCode, code int, msg string) BusError {
	return BusError{HttpCode: httpCode, ErrCode: code, Message: msg}
}

func (b BusError) GetHttpCode() int {
	if AlwaysStatusOK {
		return http.StatusOK
	}
	return b.HttpCode
}

func (b BusError) Format(params ...interface{}) BusError {
	b.Message = fmt.Sprintf(b.Message, params...)
	return b
}

func (b BusError) Wrap(err error) BusError {
	b.Errors = errors.Wrap(err, b.Message)
	return b
}

func (b BusError) Wrapf(err error, args ...interface{}) BusError {
	return b.Format(args...).Wrap(err)
}

func (b BusError) Error() string {
	return b.Message
}

func (b BusError) Stack() string {
	return fmt.Sprintf("%+v", b.Errors)
}

func (b BusError) StackAsList() []string {
	return strings.Split(b.Stack(), "\n")
}

func (b BusError) JSON(errorStack bool) JSON {
	data := JSON{
		Code:    b.ErrCode,
		Message: b.Message,
		Data:    nil,
	}

	if errorStack {
		data.ErrorStack = b.StackAsList()
	}

	return data
}
