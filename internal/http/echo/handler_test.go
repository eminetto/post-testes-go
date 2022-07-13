//go:build unit

package echo_test

import (
	"encoding/json"
	"fmt"
	"github.com/PicPay/go-test-workshop/internal/http/echo"
	"github.com/PicPay/go-test-workshop/person"
	person_mock "github.com/PicPay/go-test-workshop/person/mocks"
	weather "github.com/PicPay/go-test-workshop/weather"
	weather_mock "github.com/PicPay/go-test-workshop/weather/mocks"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	c := echo.Handlers(nil, nil, nil).NewContext(req, rec)
	c.SetPath("/hello")
	err := echo.Hello(c)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello, World!", rec.Body.String())
}

func TestGetUser(t *testing.T) {
	t.Run("status ok", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		p := []*person.Person{
			{
				ID:       1,
				Name:     "Ronnie",
				LastName: "Dio",
			},
		}
		s := person_mock.NewUseCase(t)
		s.On("Search", "dio").
			Return(p, nil).
			Once()
		c := echo.Handlers(nil, nil, nil).NewContext(req, rec)
		c.SetPath("/hello/:lastname")
		c.SetParamNames("lastname")
		c.SetParamValues("dio")
		h := echo.GetUser(s)
		err := h(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello Ronnie Dio", rec.Body.String())
	})
	t.Run("status not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		s := person_mock.NewUseCase(t)
		s.On("Search", "dio").
			Return([]*person.Person{}, nil).
			Once()
		c := echo.Handlers(nil, nil, nil).NewContext(req, rec)
		c.SetPath("/hello/:lastname")
		c.SetParamNames("lastname")
		c.SetParamValues("dio")
		h := echo.GetUser(s)
		err := h(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

}

func TestWeather(t *testing.T) {
	lat := "-27.5969"
	long := "-48.5495"
	req, _ := http.NewRequest("GET", "/weather", nil)
	t.Run("status ok", func(t *testing.T) {
		rec := httptest.NewRecorder()
		city := &weather.Weather{
			Coord: weather.Coord{
				Lon: -48.5495,
				Lat: -27.5969,
			},
			Main: weather.Main{
				Temp:      19.69,
				FeelsLike: 20.2,
				TempMin:   15.99,
				TempMax:   20.96,
				Pressure:  1013,
				Humidity:  95,
			},
			Wind: weather.Wind{
				Speed: 2.57,
				Deg:   90,
			},
			Name: "Florianópolis",
		}
		s := weather_mock.NewUseCase(t)
		s.On("Get", lat, long).
			Return(city, nil).
			Once()
		c := echo.Handlers(nil, nil, nil).NewContext(req, rec)
		c.SetPath("/weather/:lat/:long")
		c.SetParamNames("lat", "long")
		c.SetParamValues(lat, long)
		h := echo.Weather(s)
		err := h(c)
		assert.Nil(t, err)

		expected, err := json.Marshal(city)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, string(expected), rec.Body.String())
	})
	t.Run("status error", func(t *testing.T) {
		rec := httptest.NewRecorder()
		s := weather_mock.NewUseCase(t)
		s.On("Get", lat, long).
			Return(nil, fmt.Errorf("Not found")).
			Once()
		c := echo.Handlers(nil, nil, nil).NewContext(req, rec)
		c.SetPath("/weather/:lat/:long")
		c.SetParamNames("lat", "long")
		c.SetParamValues(lat, long)
		h := echo.Weather(s)
		err := h(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}
