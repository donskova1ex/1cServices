// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Swagger user management service - OpenAPI 3.0
 *
 * This is a sample some endpoints for 1C
 *
 * API version: 1.0.0
 */

package openapi

import (
	"context"
	"net/http"
)

// PDNcalculationAPIRouter defines the required methods for binding the api requests to a responses for the PDNcalculationAPI
// The PDNcalculationAPIRouter implementation should parse necessary information from the http request,
// pass the data to a PDNcalculationAPIServicer to perform the required actions, then write the service results to the http response.
type PDNcalculationAPIRouter interface {
	GetParametresByLoanId(http.ResponseWriter, *http.Request)
}

// RkoByDivisionAPIRouter defines the required methods for binding the api requests to a responses for the RkoByDivisionAPI
// The RkoByDivisionAPIRouter implementation should parse necessary information from the http request,
// pass the data to a RkoByDivisionAPIServicer to perform the required actions, then write the service results to the http response.
type RkoByDivisionAPIRouter interface {
	RkoByDivision(http.ResponseWriter, *http.Request)
}

// PDNcalculationAPIServicer defines the api actions for the PDNcalculationAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type PDNcalculationAPIServicer interface {
	GetParametresByLoanId(context.Context, string) (ImplResponse, error)
}

// RkoByDivisionAPIServicer defines the api actions for the RkoByDivisionAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type RkoByDivisionAPIServicer interface {
	RkoByDivision(context.Context, string, string) (ImplResponse, error)
}
