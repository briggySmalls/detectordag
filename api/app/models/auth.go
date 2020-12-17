package models

type Credentials struct {
	// Username for the user's account
	// required: true
	Username string `json:"username"`
	// Password for the user's account
	// required: true
	Password string `json:"password"`
}

type Token struct {
	// Token that grants access
	// required: true
	Token string `json:"token"`
	// Identifier for user's account
	// required: true
	AccountId string `json:"accountId"`
}

// swagger:parameters getAccount updateAccount getDevices registerDevice updateDevice
type TokenParameter struct {
	// A token obtained through authentication
	//
	// required: true
	// in: header
	Token string `json:"Authorization"`
}

// Credentials for logging in
// swagger:parameters auth
type AuthParameters struct {
	// Credentials for logging in
	//
	// required: true
	// in:body
	Body Credentials
}

// Successful authentication
// swagger:response tokenResponse
type TokenResponse struct {
	//in:body
	Body Token
}

// Authentication failed
// swagger:response authFailedResponse
type AuthFailedResponse struct {
	// in:body
	Body ModelError
}

// Authorization token missing or invalid
// swagger:response unauthenticatedResponse
type UnauthenticatedResponse struct {
	// in:body
	Body ModelError
}

// Token doesn't permit accessing this resource
// swagger:response unauthorizedResponse
type UnauthorisedResponse struct {
	// in:body
	Body ModelError
}
