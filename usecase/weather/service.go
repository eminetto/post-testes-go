package weather

import (
	"encoding/json"
	"fmt"
	"github.com/PicPay/go-test-workshop/entity"
	"io/ioutil"
	"net/http"
	"time"
)

type Service struct {
	client HTTPClient
	apiKey string
	url    string
}

type ServiceOption func(*Service)

func NewService(apiKey string, options ...ServiceOption) *Service {
	s := &Service{
		client: &http.Client{Timeout: time.Duration(1) * time.Second},
		apiKey: apiKey,
		url:    "https://api.openweathermap.org/data/2.5/weather?units=metric&lang=pt_br",
	}

	for _, o := range options {
		o(s)
	}
	return s
}

func WithClient(client HTTPClient) ServiceOption {
	return func(s *Service) {
		s.client = client
	}
}

func (s *Service) Get(lat, long string) (*entity.Weather, error) {
	url := fmt.Sprintf("%s&lat=%s&lon=%s&appid=%s", s.url, lat, long, s.apiKey)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var w entity.Weather
	err = json.Unmarshal(body, &w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
