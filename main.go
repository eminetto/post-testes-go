package main

import (
	"database/sql"
	"fmt"
	"github.com/PicPay/go-test-workshop/api"
	infra "github.com/PicPay/go-test-workshop/infraestructure/repository/person"
	"github.com/PicPay/go-test-workshop/usecase/person"
	"github.com/PicPay/go-test-workshop/usecase/weather"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"log"
	"os"
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
	repo := infra.NewMySQL(db)
	pService := person.NewService(repo)

	wService := weather.NewService(os.Getenv("API_KEY"))

	e := echo.New()
	e.GET("/hello", api.Hello)
	e.GET("/hello/:lastname", api.GetUser(pService))
	e.GET("/weather/:lat/:long", api.Weather(wService))
	e.Logger.Fatal(e.Start(":8000"))

}
