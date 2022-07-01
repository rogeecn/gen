## GEN: gin enhance

add some enhance feature for golang gin web framework

```golang
var (
	BusErr_UserNotFound = NewBusError(http.StatusNotFound, 1001, "User not found: %s")
	BusErr_GetParamFailed = NewBusError(http.StatusBadRequest, 1002, "get param %s failed")
)

func main() {
	svc := gin.Default()
	svc.GET("/user/:name", gen.Func(serve))
	svc.GET("/data", gen.DataFunc(serveData))
	svc.POST("/user/:id", gen.DataFunc2(serveUserUpdate,gen.Int("id", BusErr_GetParamFailed), gen.Bind(&User{})))
	svc.Run(":888")
}

func serve(ctx *gin.Context) error {
	return BusErr_UserNotFound
}

func serveData(ctx *gin.Context) (User, error) {
	return User{Name: "John"}, nil
}

func serveUserUpdate(ctx *gin.Context, id int, user *User) (User, error) {
	return user, nil
}
```