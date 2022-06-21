//go:build e2e

package api_test

import (
	"context"
	"database/sql"
	"github.com/PicPay/go-test-workshop/api"
	infra "github.com/PicPay/go-test-workshop/infraestructure/repository/person"
	usecase "github.com/PicPay/go-test-workshop/usecase/person"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserE2E(t *testing.T) {
	//fase: Configure os dados de teste
	ctx := context.Background()
	container, err := infra.SetupMysqL(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)
	db, err := sql.Open("mysql", container.URI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = infra.InitMySQL(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	defer infra.TruncateMySQL(ctx, db)

	repo := infra.NewMySQL(db)
	service := usecase.NewService(repo)
	_, err = service.Create("Ronnie", "Dio")
	assert.Nil(t, err)

	//fase: Invoque o método sendo testado
	e := echo.New()
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/hello/:lastname")
	c.SetParamNames("lastname")
	c.SetParamValues("dio")
	h := api.GetUser(service)

	//fase: Confirme que os resultados esperados são retornados
	err = h(c)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello Ronnie Dio", rec.Body.String())
}
