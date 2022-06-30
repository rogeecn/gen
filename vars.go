package gen

import "github.com/sirupsen/logrus"

var (
	AlwaysStatusOK = false
)

var LOG Logger = logrus.New()

type Logger interface {
	Error(args ...interface{})
}
