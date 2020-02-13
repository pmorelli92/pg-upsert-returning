package main

import (
	"fmt"
	"os"
	"pg-upsert-returning/postgres"
	"pg-upsert-returning/server"
)

func main() {
	customerRepo, err := postgres.NewPgCustomerRepo("postgres://postgres@127.0.0.1:5432/customer_svc?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server := server.Server{HttpAddr: "127.0.0.1:8080", CustomerRepo: customerRepo}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
