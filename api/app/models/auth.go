package models

// swagger:model credentials
type Credentials struct {

	Username string `json:"username"`

	Password string `json:"password"`
}

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

// Credentials for logging in
// swagger:parameters auth
type AuthParameters struct {
	// Credentials for logging in
	//
	// in:body
	Body Credentials
}

// Successful authentication
// swagger:response tokenResponse
type TokenResponse struct {
	//in:body
	Body Token
}

// Authorization token missing or invalid
// swagger:response unauthenticatedResponse
type AuthFailedResponse struct {
	// in:body
	Body ModelError
}

// Token doesn't permit accessing this resource
// swagger:response unauthorizedResponse
type UnauthorisedResponse struct {
	// in:body
	Body ModelError
}
