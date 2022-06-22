//go:build unit

package api_test

import (
	"encoding/json"
	"fmt"
	"github.com/PicPay/go-test-workshop/api"
	"github.com/PicPay/go-test-workshop/entity"
	person "github.com/PicPay/go-test-workshop/usecase/person/mocks"
	weather "github.com/PicPay/go-test-workshop/usecase/weather/mocks"
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
	t.Run("status ok", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		p := []*entity.Person{
			{
				ID:       1,
				Name:     "Ronnie",
				LastName: "Dio",
			},
		}
		s := person.NewUseCase(t)
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
	})
	t.Run("status not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		s := person.NewUseCase(t)
		s.On("Search", "dio").
			Return([]*entity.Person{}, nil).
			Once()
		c := e.NewContext(req, rec)
		c.SetPath("/hello/:lastname")
		c.SetParamNames("lastname")
		c.SetParamValues("dio")
		h := api.GetUser(s)
		err := h(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

}

func TestWeather(t *testing.T) {
	e := echo.New()
	lat := "-27.5969"
	long := "-48.5495"
	req, _ := http.NewRequest("GET", "/weather", nil)
	t.Run("status ok", func(t *testing.T) {
		rec := httptest.NewRecorder()
		city := &entity.Weather{
			Coord: entity.Coord{
				Lon: -48.5495,
				Lat: -27.5969,
			},
			Main: entity.Main{
				Temp:      19.69,
				FeelsLike: 20.2,
				TempMin:   15.99,
				TempMax:   20.96,
				Pressure:  1013,
				Humidity:  95,
			},
			Wind: entity.Wind{
				Speed: 2.57,
				Deg:   90,
			},
			Name: "Florianópolis",
		}
		s := weather.NewUseCase(t)
		s.On("Get", lat, long).
			Return(city, nil).
			Once()
		c := e.NewContext(req, rec)
		c.SetPath("/weather/:lat/:long")
		c.SetParamNames("lat", "long")
		c.SetParamValues(lat, long)
		h := api.Weather(s)
		err := h(c)
		assert.Nil(t, err)

		expected, err := json.Marshal(city)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, string(expected), rec.Body.String())
	})
	t.Run("status error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		s := weather.NewUseCase(t)
		s.On("Get", lat, long).
			Return(nil, fmt.Errorf("Not found")).
			Once()
		c := e.NewContext(req, rec)
		c.SetPath("/weather/:lat/:long")
		c.SetParamNames("lat", "long")
		c.SetParamValues(lat, long)
		h := api.Weather(s)
		err := h(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}
