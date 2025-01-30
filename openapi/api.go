package openapi

import (
	"context"
	"net/http"
)

// PDNcalculationAPIRouter defines the required methods for binding the 1c_api requests to a responses for the PDNcalculationAPI
// The PDNcalculationAPIRouter implementation should parse necessary information from the http request,
// pass the data to a PDNcalculationAPIServicer to perform the required actions, then write the service results to the http response.
type PDNcalculationAPIRouter interface {
	GetParametresByLoanId(http.ResponseWriter, *http.Request)
}

// PDNcalculationAPIServicer defines the 1c_api actions for the PDNcalculationAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type PDNcalculationAPIServicer interface {
	GetParametresByLoanId(context.Context, string) (ImplResponse, error)
}
