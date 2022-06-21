//go:build unit

package api_test

import (
	"github.com/PicPay/go-test-workshop/api"
	"github.com/PicPay/go-test-workshop/entity"
	"github.com/PicPay/go-test-workshop/usecase/person/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	e := echo.New()
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/hello")
	err := api.Hello(c)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello, World!", rec.Body.String())
}

func TestGetUser(t *testing.T) {
	e := echo.New()
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	p := []*entity.Person{
		{
			ID:       1,
			Name:     "Ronnie",
			LastName: "Dio",
		},
	}
	s := mocks.NewUseCase(t)
	s.On("Search", "dio").
		Return(p, nil).
		Once()
	c := e.NewContext(req, rec)
	c.SetPath("/hello/:lastname")
	c.SetParamNames("lastname")
	c.SetParamValues("dio")
	h := api.GetUser(s)
	err := h(c)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello Ronnie Dio", rec.Body.String())
}
