package main

import (
	"log"
	"net/http"

	openapi "github.com/donskova1ex/1cServices/openapi"
)

func main() {
	log.Printf("Server started")

	PDNcalculationAPIService := openapi.NewPDNcalculationAPIService()
	PDNcalculationAPIController := openapi.NewPDNcalculationAPIController(PDNcalculationAPIService)

	router := openapi.NewRouter(PDNcalculationAPIController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
