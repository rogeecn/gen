package gen

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Int(key string, errRet BusError) func(*gin.Context) (int, error) {
	return func(ctx *gin.Context) (int, error) {
		v, ok := ctx.Params.Get(key)
		if !ok {
			return 0, errRet.Format(key)
		}

		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, errRet.Format(key).Wrap(err)
		}
		return i, nil
	}
}

func String(key string, errRet BusError) func(*gin.Context) (string, error) {
	return func(ctx *gin.Context) (string, error) {
		v, ok := ctx.Params.Get(key)
		if !ok {
			return "", errRet.Format(key)
		}

		return v, nil
	}
}

func Bind[T any](param T, errRet BusError) func(*gin.Context) (T, error) {
	return func(ctx *gin.Context) (T, error) {
		err := ctx.ShouldBind(&param)
		if err != nil {
			return param, errRet.Wrap(err)

		}

		return param, nil
	}
}
