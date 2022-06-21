package weather

import (
	"github.com/PicPay/go-test-workshop/entity"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type UseCase interface {
	Get(lat, long string) (*entity.Weather, error)
}
