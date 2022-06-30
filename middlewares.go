package gen

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrProc  func(ctx *gin.Context, err BusError)
	DataProc func(ctx *gin.Context, data interface{})
)

func init() {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
}

func Error(f func(*gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := f(ctx)
		if err != nil {
			ErrProc(ctx, err.(BusError))
			return
		}

		ctx.Next()
	}
}

func DataError(f func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		if err != nil {
			ErrProc(ctx, err.(BusError))
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}

func defaultErrProc(ctx *gin.Context, err BusError) {
	if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		LOG.Error(err.Stack())
		ctx.JSON(err.GetHttpCode(), err.JSON())
		return
	}

	ctx.AbortWithError(err.GetHttpCode(), err)
}

func defaultDataProc(ctx *gin.Context, data interface{}) {
	ctx.JSONP(http.StatusOK, JSON{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	})
}
