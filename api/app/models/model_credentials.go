package models

// swagger:model credentials
type Credentials struct {

	Username string `json:"username"`

	Password string `json:"password"`
}

// Credentials for logging in
// swagger:parameters auth
type AuthParameters struct {
	// Credentials for logging in
	//
	// in:body
	Body Credentials
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
