package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/PicPay/go-test-workshop/internal/api"
	"github.com/PicPay/go-test-workshop/internal/http/echo"
	"github.com/PicPay/go-test-workshop/person"
	"github.com/PicPay/go-test-workshop/person/mysql"
	"github.com/PicPay/go-test-workshop/weather"
	logger "github.com/PicPay/lib-go-logger"
	_ "github.com/go-sql-driver/mysql"
)

const ( //@todo pegar essa informação de variáveis de ambiente
	dbUser         = "workshop"
	dbPassword     = "workshop"
	database       = "workshop"
	dbRootPassword = "db-root-password"
)

func main() {
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, "localhost", "3306", database)
	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		log.Fatal(err)
	}
	repo := mysql.NewMySQL(db)
	pService := person.NewService(repo)

	wService := weather.NewService(os.Getenv("API_KEY"))

	l := logger.New()
	h := echo.Handlers(l, pService, wService)
	err = api.Start(l, "8000", h)
	if err != nil {
		l.Fatal("error running api", err)
	}
}
