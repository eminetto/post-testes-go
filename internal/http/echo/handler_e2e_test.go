//go:build e2e

package echo_test

import (
	"context"
	"database/sql"
	"github.com/PicPay/go-test-workshop/internal/http/echo"
	"github.com/PicPay/go-test-workshop/person"
	"github.com/PicPay/go-test-workshop/person/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserE2E(t *testing.T) {
	//fase: Configure os dados de teste
	ctx := context.Background()
	container, err := person.SetupMysqL(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)
	db, err := sql.Open("mysql", container.URI)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	err = person.InitMySQL(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	defer person.TruncateMySQL(ctx, db)

	repo := mysql.NewMySQL(db)
	service := person.NewService(repo)
	_, err = service.Create("Ronnie", "Dio")
	assert.Nil(t, err)

	//fase: Invoque o método sendo testado
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	c := echo.Handlers(nil, nil, nil).NewContext(req, rec)
	c.SetPath("/hello/:lastname")
	c.SetParamNames("lastname")
	c.SetParamValues("dio")
	h := echo.GetUser(service)

	//fase: Confirme que os resultados esperados são retornados
	err = h(c)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello Ronnie Dio", rec.Body.String())
}
