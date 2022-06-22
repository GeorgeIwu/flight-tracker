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
	"service-catalog/service"
)


func main() {
	e := echo.New()
	client, err := ent.Open("mysql", "kong:password@tcp(localhost)/catalog?parseTime=True")
	if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
	}

	timeoutContext := time.Duration(60) * time.Second
	uc := service.NewServiceUsecase(client, timeoutContext)
	service.NewServiceHandler(e, uc)


	fmt.Printf("Running app on port 8000")
	log.Fatal(http.ListenAndServe(":8000", e))
}
