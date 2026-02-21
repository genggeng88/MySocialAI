package main

import (
	"fmt"
	"log"
	"mysocialai/backend"
	"mysocialai/handler"
	"net/http"
)

func main() {
	fmt.Println("started-service")

	backend.InitElasticsearchBackend()

	log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}