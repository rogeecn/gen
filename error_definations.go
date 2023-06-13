package gen

import "net/http"

var (
	PathParamError   = NewBusError(http.StatusBadRequest, 1001, "path param error %s")
	QueryParamError  = NewBusError(http.StatusBadRequest, 1002, "query param error %s")
	BodyParamError   = NewBusError(http.StatusBadRequest, 1003, "body param error %s")
	HeaderParamError = NewBusError(http.StatusBadRequest, 1004, "header param error %s")
)
