package gen

import (
	"net/http"
)

var (
	PathParamError   = NewBusError(http.StatusBadRequest, http.StatusBadRequest, "path param error")
	QueryParamError  = NewBusError(http.StatusBadRequest, http.StatusBadRequest, "query param error")
	BodyParamError   = NewBusError(http.StatusBadRequest, http.StatusBadRequest, "body param error")
	HeaderParamError = NewBusError(http.StatusBadRequest, http.StatusBadRequest, "header param error")

	StatusNotFoundErr = NewBusError(http.StatusNotFound, http.StatusNotFound, "resource not found")
)

func defaultErrorConvert(err error) BusError {
	return NewBusError(500, 500, err.Error()).Wrap(err)
}
