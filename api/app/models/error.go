package models

type ModelError struct {
	// Description of the error
	// required: true
	Error_ string `json:"error"`
}
