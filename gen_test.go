package gen

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type respUser struct {
	ID   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type testApi struct {
}

var (
	Err_Test = NewBusError(http.StatusBadRequest, 1001, "TestErr")
)

func (t *testApi) Func(ctx *gin.Context) error {
	return Err_Test.Wrap(errors.New("TestStack"))
}

func (t *testApi) FuncP1(ctx *gin.Context, uid int) error {
	return Err_Test.Wrap(errors.New("TestStack"))
}

func (t *testApi) Data(ctx *gin.Context) (*respUser, error) {
	return &respUser{Name: "TestName"}, nil
}

func (t *testApi) DataP1(ctx *gin.Context, uid int) (*respUser, error) {
	return &respUser{ID: uid, Name: "TestName"}, nil
}

func (t *testApi) DataP2(ctx *gin.Context, uid int, name string) (*respUser, error) {
	return &respUser{ID: uid, Name: name}, nil
}

func (t *testApi) DataP1Form(ctx *gin.Context, user *respUser) (*respUser, error) {
	return user, nil
}

func Test_Func(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	Func(api.Func)(ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp JSON
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1001, resp.Code)
	assert.Equal(t, "TestErr", resp.Message)
	assert.Nil(t, resp.Data)
	assert.NotNil(t, resp.ErrorStack)
}

func Test_Func_P1(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	ErrParamNotExist := NewBusError(http.StatusBadRequest, 10001, "%d 参数不存在")

	svc := gin.Default()
	svc.GET("/step1/:uid", Func1(api.FuncP1, Int("uidd", ErrParamNotExist)))
	svc.GET("/step2/:uid", Func1(api.FuncP1, Int("uid", ErrParamNotExist)))

	req := httptest.NewRequest(http.MethodGet, "/step1/100", nil)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req = httptest.NewRequest(http.MethodGet, "/step2/100", nil)
	w = httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_DataFunc(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	DataFunc(api.Data)(ctx)
	assert.Equal(t, http.StatusOK, rec.Code)

	t.Logf("resp body: %s\n\n", rec.Body.String())

	var resp JSON
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	b, err := json.Marshal(resp.Data)
	assert.NoError(t, err)

	var data respUser
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err)
	assert.Equal(t, "TestName", data.Name)
}

func Test_DataFunc_P1(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	ErrParamNotExist := NewBusError(http.StatusBadRequest, 10001, "%s 参数不存在")

	svc := gin.Default()
	svc.GET("/step1/:uid", DataFunc1(api.DataP1, Int("uidd", ErrParamNotExist)))
	svc.GET("/step2/:uid", DataFunc1(api.DataP1, Int("uid", ErrParamNotExist)))

	req := httptest.NewRequest(http.MethodGet, "/step1/100", nil)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	t.Logf("resp body: %s\n\n", w.Body.String())

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// step2 test
	req = httptest.NewRequest(http.MethodGet, "/step2/100", nil)
	w = httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("resp body: %s\n\n", w.Body.String())

	var resp JSON
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	b, err := json.Marshal(resp.Data)
	assert.NoError(t, err)

	var data respUser
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err)
	assert.Equal(t, 100, data.ID)
}

func Test_DataFunc_P2(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	ErrParamNotExist := NewBusError(http.StatusBadRequest, 10001, "%s 参数不存在")

	svc := gin.Default()
	svc.GET("/:uid/:name", DataFunc2(api.DataP2, Int("uid", ErrParamNotExist), String("name", ErrParamNotExist)))

	req := httptest.NewRequest(http.MethodGet, "/100/ZhangSan", nil)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	t.Logf("resp body: %s\n\n", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

	var resp JSON
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	b, err := json.Marshal(resp.Data)
	assert.NoError(t, err)

	var data respUser
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err)
	assert.Equal(t, 100, data.ID)
	assert.Equal(t, "ZhangSan", data.Name)
}

func Test_DataFunc_P1_POST(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	ErrBindParam := NewBusError(http.StatusBadRequest, 10001, "参数绑定失败")

	svc := gin.Default()
	svc.POST("/user", DataFunc1(api.DataP1Form, BindBody(&respUser{}, ErrBindParam)))
	svc.POST("/user2", DataFunc1(api.DataP1Form, BindBody(&respUser{}, ErrBindParam)))

	form := &url.Values{}
	form.Add("name", "TestName")
	form.Add("id", "100")
	t.Logf("body: %s", form.Encode())
	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	t.Logf("resp body: %s\n\n", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
	var resp JSON
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	b, err := json.Marshal(resp.Data)
	assert.NoError(t, err)

	var data respUser
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err)
	assert.Equal(t, 100, data.ID)
	assert.Equal(t, "TestName", data.Name)

	// step2 test
	b, err = json.Marshal(respUser{ID: 100, Name: "TestName"})
	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/user2", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("resp body: %s\n\n", w.Body.String())

	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)

	b, err = json.Marshal(resp.Data)
	assert.NoError(t, err)

	err = json.Unmarshal(b, &data)
	assert.NoError(t, err)
	assert.Equal(t, 100, data.ID)
}
