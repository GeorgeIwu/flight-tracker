package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"service-catalog/ent"
	"service-catalog/middleware"
	"service-catalog/service"
)

//move to environment variable
const (
	dbHost     = "localhost"
	dbUser     = "kong"
	dbPassword = "password"
	dbType    = "mysql"
	dbProtocol = "tcp"
	dbName = "catalog"
)

func main() {
	sqlInfo := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=True", dbUser, dbPassword, dbProtocol, dbHost, dbName)
	client, err := ent.Open(dbType, sqlInfo)
	if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
	}

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.Auth)
	timeoutContext := time.Duration(60) * time.Second
	repo := service.NewServiceRepo(client, timeoutContext)
	usecase := service.NewServiceUsecase(repo)
	service.NewServiceHandler(e, usecase)


	fmt.Printf("Running app on port 8000")
	log.Fatal(http.ListenAndServe(":8000", e))
}
