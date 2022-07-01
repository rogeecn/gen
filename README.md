## GEN: gin enhance

add some enhance feature for golang gin web framework

```golang
var (
	BusErr_UserNotFound = NewBusError(http.StatusNotFound, 1001, "User not found: %s")
)

func main() {
	svc := gin.Default()
	svc.GET("/user/:name", gen.Func(serve))
	svc.GET("/data", gen.DataFunc(serveData))
	svc.Run(":888")
}

func serve(ctx *gin.Context) error {
	return BusErr_UserNotFound
}

func serveData(ctx *gin.Context) (User, error) {
	return User{Name: "John"}, nil
}
```