package main

import (
	"fmt"
	"os"
	"net/http"

	"axiata_test/routes"

	"github.com/joho/godotenv"
)

func main () {
	if err := godotenv.Load(); err != nil {
		fmt.Println(" No .env file is found.. ")
	}
	routes.Register()

	http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), nil)
}