package weather_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/PicPay/go-test-workshop/weather"
	"github.com/PicPay/go-test-workshop/weather/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	client := mocks.NewHTTPClient(t)
	lat := "-48.5495"
	long := "-27.5969"
	url := "https://api.openweathermap.org/data/2.5/weather?units=metric&lang=pt_br"
	apiKey := "fake"

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s&lat=%s&lon=%s&appid=%s", url, lat, long, apiKey), nil)
	assert.Nil(t, err)
	json := `{"coord":{"lon":-48.5495,"lat":-27.5969},"weather":[{"id":211,"main":"Thunderstorm","description":"trovoadas","icon":"11d"}],"base":"stations","main":{"temp":19.69,"feels_like":20.2,"temp_min":15.99,"temp_max":20.96,"pressure":1013,"humidity":95},"visibility":10000,"wind":{"speed":2.57,"deg":90},"clouds":{"all":75},"dt":1655836456,"sys":{"type":2,"id":2018322,"country":"BR","sunrise":1655805850,"sunset":1655843264},"timezone":-10800,"id":3463237,"name":"Florianópolis","cod":200}`
	body := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	client.On("Do", request).
		Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil).
		Once()
	s := weather.NewService(apiKey,
		weather.WithClient(client),
	)
	expected := &weather.Weather{
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
	w, err := s.Get(lat, long)
	assert.Nil(t, err)
	assert.Equal(t, expected, w)
}
