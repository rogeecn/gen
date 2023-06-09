package gen

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/constraints"
)

func Integer[T constraints.Integer](key string, errRet BusError) func(*gin.Context) (T, error) {
	return func(ctx *gin.Context) (T, error) {
		v, ok := ctx.Params.Get(key)
		if !ok {
			return 0, errRet.Format(key)
		}

		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, errRet.Format(key).Wrap(err)
		}
		return T(i), nil
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

func Body[T any](param T, errRet BusError) func(*gin.Context) (T, error) {
	return func(ctx *gin.Context) (T, error) {
		if contentType := ctx.Request.Header.Get("Content-Type"); strings.TrimSpace(contentType) == "" {
			ctx.Request.Header.Set("Content-Type", DefaultBindBodyMIME)
		}

		err := ctx.ShouldBind(&param)
		if err != nil {
			return param, errRet.Wrap(err)
		}

		return param, nil
	}
}

func Query[T any](param T, errRet BusError) func(*gin.Context) (T, error) {
	return func(ctx *gin.Context) (T, error) {
		err := ctx.ShouldBindQuery(&param)
		if err != nil {
			return param, errRet.Wrap(err)
		}

		return param, nil
	}
}

func Header[T any](param T, errRet BusError) func(*gin.Context) (T, error) {
	return func(ctx *gin.Context) (T, error) {
		err := ctx.ShouldBindHeader(&param)
		if err != nil {
			return param, errRet.Wrap(err)
		}

		return param, nil
	}
}
