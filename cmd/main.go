package main

import (
	"fmt"
	"os"
	"pg-upsert-returning/postgres"
	"pg-upsert-returning/server"
)

func main() {
	customerRepo, err := postgres.NewPgCustomerRepo("postgres://postgres@db:5432/customer_svc?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server := server.Server{HttpAddr: ":8080", CustomerRepo: customerRepo}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
