package main

import (
	"database/sql"
	"fmt"
	"github.com/PicPay/go-test-workshop/api"
	infra "github.com/PicPay/go-test-workshop/infraestructure/repository/person"
	usecase "github.com/PicPay/go-test-workshop/usecase/person"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"log"
)

const ( //@todo pegar essa informação de variáveis de ambiente
	dbUser         = "workshop"
	dbPassword     = "workshop"
	database       = "workshop"
	dbRootPassword = "db-root-password"
)

func main() {
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", "root", dbRootPassword, "localhost", "3306", database)
	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		log.Fatal(err)
	}
	repo := infra.NewMySQL(db)
	service := usecase.NewService(repo)
	e := echo.New()
	e.GET("/hello", api.Hello)
	e.GET("/hello/:lastname", api.GetUser(service))
	e.Logger.Fatal(e.Start(":8000"))

}
