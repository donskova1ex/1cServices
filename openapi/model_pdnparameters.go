// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Swagger user management service - OpenAPI 3.0
 *
 * This is a sample some endpoints for 1C
 *
 * API version: 1.0.0
 */

package openapi

type Pdnparameters struct {
	LoanId string `json:"LoanId"`

	Incomes float32 `json:"Incomes"`

	Expenses float32 `json:"Expenses"`

	IncomesTypeId string `json:"IncomesTypeId"`

	AverageRegionIncomes float32 `json:"AverageRegionIncomes"`
}

// AssertPdnparametersRequired checks if the required fields are not zero-ed
func AssertPdnparametersRequired(obj Pdnparameters) error {
	return nil
}

// AssertPdnparametersConstraints checks if the values respects the defined constraints
func AssertPdnparametersConstraints(obj Pdnparameters) error {
	return nil
}
