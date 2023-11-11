package dto

// UserCreateRequest is the request body for the user create endpoint.
type UserCreateRequest struct {
	Email    string `json:"email"          validate:"required,email"`
	Username string `json:"username"       validate:"required,min=3,max=21"`
	Password string `json:"password"       validate:"required,min=6,max=55"`
}
