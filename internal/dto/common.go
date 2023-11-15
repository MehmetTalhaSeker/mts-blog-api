package dto

type RequestWithID struct {
	ID string `param:"id" validate:"required"`
}

type ResponseWithID struct {
	ID string `json:"id"`
}
