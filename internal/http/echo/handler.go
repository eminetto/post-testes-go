package echo

import (
	"fmt"
	"net/http"

	"github.com/PicPay/go-test-workshop/person"
	"github.com/PicPay/go-test-workshop/weather"
	logger "github.com/PicPay/lib-go-logger"
	"github.com/labstack/echo/v4"
)

func Handlers(l *logger.Logger, pService person.UseCase, wService weather.UseCase) *echo.Echo {
	e := echo.New()
	e.GET("/hello", Hello)
	e.GET("/hello/:lastname", GetUser(pService))
	e.GET("/weather/:lat/:long", Weather(wService))
	return e
}

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func GetUser(s person.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		lastname := c.Param("lastname")
		people, err := s.Search(lastname)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if len(people) == 0 {
			return c.String(http.StatusNotFound, "not found")
		}
		return c.String(http.StatusOK, fmt.Sprintf("Hello %s %s", people[0].Name, people[0].LastName))
	}
}

func Weather(s weather.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		lat := c.Param("lat")
		long := c.Param("long")
		w, err := s.Get(lat, long)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, w)
	}
}
