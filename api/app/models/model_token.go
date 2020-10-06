package models

// swagger:model token
type Token struct {

	Token string `json:"token"`

	AccountId string `json:"accountId"`
}

// swagger:parameters getAccount updateAccount getDevices registerDevice updateDevice
type TokenParameter struct {
	// A token obtained through authentication
	//
	// required: true
	// in: header
	Token string `json:"token"`
}
