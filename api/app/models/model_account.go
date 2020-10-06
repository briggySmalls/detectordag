package models

type Emails struct {
	// The emails associated with the account
	Emails []string `json:"emails"`
}

// swagger:model account
type Account struct {
	// The username of the account
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

// swagger:parameters getAccount getDevices updateAccount registerDevice
type AccountParameter struct {
	// ID of account that is to be queried
	//
	// required: true
	// in: path
	AccountID string `json:"accountId"`
}
