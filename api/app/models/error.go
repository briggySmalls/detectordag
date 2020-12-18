package models

type ModelError struct {
	// Description of the error
	// required: true
	// example: Something went terribly wrong
	Error_ string `json:"error"`
}
