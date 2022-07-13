package api

import (
	"net/http"
	"testing"

	logger "github.com/PicPay/lib-go-logger"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct{}

func (m mockHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestStart(t *testing.T) {
	t.Run("retorna erro ao executar ListenAndServe com porta inválida", func(t *testing.T) {
		logger := logger.New()
		err := Start(logger, "abacate", mockHandler{})
		assert.Contains(t, err.Error(), "listen tcp: lookup tcp/abacate: nodename nor servname provided, or not known")
	})
}
