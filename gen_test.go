package gen

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_Err(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("X-Requested-With", "XMLHttpRequest")
	Func(api.Err)(ctx)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	t.Logf("resp body: %s\n\n", rec.Body.String())

	var resp JSON
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1001, resp.Code)
	assert.Equal(t, "TestErr", resp.Message)
	assert.Nil(t, resp.Data)
}

func Test_GenericDataErr(t *testing.T) {
	ErrProc = defaultErrProc
	DataProc = defaultDataProc
	api := &testApi{}

	rec := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(rec)

	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("X-Requested-With", "XMLHttpRequest")
	DataFunc(api.DataStruct)(ctx)
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

type respUser struct {
	Name string `json:"name"`
}

type testApi struct {
}

var (
	Err_Test = NewBusError(http.StatusBadRequest, 1001, "TestErr")
)

func (t *testApi) Err(ctx *gin.Context) error {
	return Err_Test.Wrap(errors.New("TestStack"))
}

func (t *testApi) DataStruct(ctx *gin.Context) (*respUser, error) {
	return &respUser{Name: "TestName"}, nil
}
