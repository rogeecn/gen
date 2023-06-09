package gen

import "net/http"

var (
	PathParamError   = NewBusError(http.StatusBadRequest, 1001, "path param error")
	QueryParamError  = NewBusError(http.StatusBadRequest, 1002, "query param error")
	BodyParamError   = NewBusError(http.StatusBadRequest, 1003, "body param error")
	HeaderParamError = NewBusError(http.StatusBadRequest, 1004, "header param error")
)
