package models

type Emails struct {
	// The emails associated with the account
	// required: true
	// example: ["jane@example.com", "john@example.com"]
	Emails []string `json:"emails"`
}

type Account struct {
	// The username of the account
	// required: true
	// example: user@example.com
	Username string `json:"username"`
	Emails
}

// Successful account retrieval
// swagger:response getAccountResponse
type AccountResponse struct {
	// in: body
	Body Account
}

// Account with that ID not found
// swagger:response accountNotFoundResponse
type AccountNotFoundResponse struct {
	// in: body
	Body ModelError
}

// swagger:parameters getAccount getDevices updateAccount
type AccountParameter struct {
	// ID of account that is to be queried
	//
	// required: true
	// in: path
	AccountID string `json:"accountId"`
}

// swagger:parameters updateAccount
type EmailsParameter struct {
	// Properties to update about the account
	//
	// required: true
	// in: body
	Body Emails
}
