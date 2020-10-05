package models

// swagger:model account
type Account struct {
	// The username of the account
	Username string `json:"username"`
	// The emails associated with the account
	Emails []string `json:"emails"`
}
