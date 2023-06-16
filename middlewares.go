package gen

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrProc          func(ctx *gin.Context, err error)
	DataProc         func(ctx *gin.Context, data interface{})
	ErrorConvertProc func(error) BusError
)

func init() {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	ErrorConvertProc = defaultErrorConvert
}

func defaultErrProc(ctx *gin.Context, err error) {
	var busErr BusError
	switch err.(type) {
	case BusError:
		busErr = err.(BusError)
	default:
		busErr = ErrorConvertProc(err)
	}

	_, _ = gin.DefaultErrorWriter.Write([]byte(busErr.Stack()))
	ctx.JSON(busErr.GetHttpCode(), busErr.JSON(ctx, gin.IsDebugging()))
}

func defaultDataProc(ctx *gin.Context, data interface{}) {
	// var json JsonResponse
	// json = &JSON{}

	// if v, ok := ctx.Get(JsonResponseKey); ok && v != nil {
	// 	json = v.(JsonResponse)
	// }

	// ctx.JSONP(http.StatusOK, json.SetCode(http.StatusOK).SetMessage("OK").SetData(data))
	ctx.JSONP(http.StatusOK, data)
}

func Func(f func(*gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := f(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		DataProc(ctx, nil)
		ctx.Next()
	}
}

func Func1[P1 any](
	f func(*gin.Context, P1) error,
	pf1 func(*gin.Context) (P1, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		err = f(ctx, p)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		DataProc(ctx, nil)
		ctx.Next()
	}
}

func Func2[P1 any, P2 any](
	f func(*gin.Context, P1, P2) error,
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		err = f(ctx, p1, p2)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		DataProc(ctx, nil)
		ctx.Next()
	}
}

func Func3[P1 any, P2 any, P3 any](
	f func(*gin.Context, P1, P2, P3) error,
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		err = f(ctx, p1, p2, p3)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		DataProc(ctx, nil)
		ctx.Next()
	}
}

func Func4[P1 any, P2 any, P3 any, P4 any](
	f func(*gin.Context, P1, P2, P3, P4) error,
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
	pf4 func(*gin.Context) (P4, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p4, err := pf4(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		err = f(ctx, p1, p2, p3, p4)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		DataProc(ctx, nil)
		ctx.Next()
	}
}

func Func5[P1 any, P2 any, P3 any, P4 any, P5 any](
	f func(*gin.Context, P1, P2, P3, P4, P5) error,
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
	pf4 func(*gin.Context) (P4, error),
	pf5 func(*gin.Context) (P5, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p4, err := pf4(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p5, err := pf5(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		err = f(ctx, p1, p2, p3, p4, p5)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		DataProc(ctx, nil)
		ctx.Next()
	}
}

func Func6[P1 any, P2 any, P3 any, P4 any, P5 any, P6 any](
	f func(*gin.Context, P1, P2, P3, P4, P5, P6) error,
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
	pf4 func(*gin.Context) (P4, error),
	pf5 func(*gin.Context) (P5, error),
	pf6 func(*gin.Context) (P6, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p4, err := pf4(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p5, err := pf5(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p6, err := pf6(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		err = f(ctx, p1, p2, p3, p4, p5, p6)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		DataProc(ctx, nil)
		ctx.Next()
	}
}

func DataFunc[T any](f func(*gin.Context) (T, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}

func DataFunc1[T any, P1 any](f func(*gin.Context, P1) (T, error), pf1 func(*gin.Context) (P1, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		data, err := f(ctx, p)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}

func DataFunc2[T any, P1 any, P2 any](
	f func(*gin.Context, P1, P2) (T, error),
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		data, err := f(ctx, p1, p2)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}

func DataFunc3[T any, P1 any, P2 any, P3 any](
	f func(*gin.Context, P1, P2, P3) (T, error),
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		data, err := f(ctx, p1, p2, p3)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}

func DataFunc4[T any, P1 any, P2 any, P3 any, P4 any](
	f func(*gin.Context, P1, P2, P3, P4) (T, error),
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
	pf4 func(*gin.Context) (P4, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p4, err := pf4(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		data, err := f(ctx, p1, p2, p3, p4)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}

func DataFunc5[T any, P1 any, P2 any, P3 any, P4 any, P5 any](
	f func(*gin.Context, P1, P2, P3, P4, P5) (T, error),
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
	pf4 func(*gin.Context) (P4, error),
	pf5 func(*gin.Context) (P5, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p4, err := pf4(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p5, err := pf5(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		data, err := f(ctx, p1, p2, p3, p4, p5)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}

func DataFunc6[T any, P1 any, P2 any, P3 any, P4 any, P5 any, P6 any](
	f func(*gin.Context, P1, P2, P3, P4, P5, P6) (T, error),
	pf1 func(*gin.Context) (P1, error),
	pf2 func(*gin.Context) (P2, error),
	pf3 func(*gin.Context) (P3, error),
	pf4 func(*gin.Context) (P4, error),
	pf5 func(*gin.Context) (P5, error),
	pf6 func(*gin.Context) (P6, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		p1, err := pf1(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p2, err := pf2(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		p3, err := pf3(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p4, err := pf4(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p5, err := pf5(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		p6, err := pf6(ctx)
		if err != nil {
			ErrProc(ctx, err)
			return
		}
		data, err := f(ctx, p1, p2, p3, p4, p5, p6)
		if err != nil {
			ErrProc(ctx, err)
			return
		}

		DataProc(ctx, data)
		ctx.Next()
	}
}
