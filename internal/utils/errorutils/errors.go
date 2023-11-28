package errorutils

import (
	"errors"
	"fmt"
	"net/http"
)

// Auth Errors.
var (
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrEmailRequired        = errors.New("email is required")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("expired token")
	ErrLongPassword         = errors.New("password should be less then 55")
	ErrLongUsername         = errors.New("username  should be less then 21")
	ErrLoginFailed          = errors.New("email or password is incorrect")
	ErrMissingAuthHeader    = errors.New("missing authorization header")
	ErrShortPassword        = errors.New("password too short")
	ErrShortUsername        = errors.New("username too short")
	ErrEmailNotFound        = errors.New("email not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrUsernameRequired     = errors.New("username is required")
)

// Common Errors.
var (
	ErrBadRequest          = errors.New("bad request")
	ErrBinding             = errors.New("something went wrong during data binding")
	ErrEmptyID             = errors.New("ID can't be empty")
	ErrInvalidID           = errors.New("invalid ID")
	ErrJSONDecode          = errors.New("json decode error")
	ErrJSONEncode          = errors.New("json encode error")
	ErrJSONMarshal         = errors.New("json marshal error")
	ErrJSONUnmarshal       = errors.New("json unmarshal error")
	ErrLongPaginationSize  = errors.New("size should be less than 100")
	ErrShortPaginationSize = errors.New("size should be more than 1")
	ErrUnexpected          = errors.New("unexpected error")
)

// User Errors.
var (
	ErrUserCount    = errors.New("user count failed")
	ErrUserCreate   = errors.New("user create failed")
	ErrUserDelete   = errors.New("user delete failed")
	ErrUserRead     = errors.New("user read failed")
	ErrUserReads    = errors.New("user reads failed")
	ErrUserSearch   = errors.New("user search failed")
	ErrUserUpdate   = errors.New("user update failed")
	ErrUserNotFound = errors.New("user not found")
)

// Post Errors.
var (
	ErrPostCount    = errors.New("post count failed")
	ErrPostCreate   = errors.New("post create failed")
	ErrPostDelete   = errors.New("post delete failed")
	ErrPostRead     = errors.New("post read failed")
	ErrPostReads    = errors.New("post reads failed")
	ErrPostUpdate   = errors.New("post update failed")
	ErrPostNotFound = errors.New("post not found")
)

// Unorganized Errors.
var (
	ErrFailedRead        = errors.New("we couldn't read your request. Please try again")
	ErrFailedSave        = errors.New("we couldn't save your request. Please try again")
	ErrInvalidPassword   = errors.New("wrong password")
	ErrInvalidQueryParam = errors.New("your requested query params are invalid. Please try again")
	ErrInvalidRequest    = errors.New("invalid request")
	ErrInvalidRole       = errors.New("invalid role")
	ErrUnauthorized      = errors.New("unauthorized user")
)

// errorCodes a map to store error codes.
var errorCodes = map[error]string{
	// Auth
	ErrEmailAlreadyTaken:    ErrCodeEmailAlreadyTaken,
	ErrEmailRequired:        ErrCodeEmailRequired,
	ErrInvalidToken:         ErrCodeInvalidToken,
	ErrExpiredToken:         ErrCodeExpiredToken,
	ErrLongPassword:         ErrCodeLongPassword,
	ErrLongUsername:         ErrCodeLongUsername,
	ErrLoginFailed:          ErrCodeLoginFailed,
	ErrMissingAuthHeader:    ErrCodeMissingAuthHeader,
	ErrShortPassword:        ErrCodeShortPassword,
	ErrShortUsername:        ErrCodeShortUsername,
	ErrUserNotFound:         ErrCodeUserNotFound,
	ErrEmailNotFound:        ErrCodeEmailNotFound,
	ErrUsernameAlreadyTaken: ErrCodeUsernameAlreadyTaken,
	ErrUsernameRequired:     ErrCodeUsernameRequired,

	// Common
	ErrBadRequest:          ErrCodeBadRequest,
	ErrBinding:             ErrCodeBinding,
	ErrEmptyID:             ErrCodeEmptyID,
	ErrInvalidID:           ErrCodeInvalidID,
	ErrJSONDecode:          ErrCodeJSONDecode,
	ErrJSONEncode:          ErrCodeJSONEncode,
	ErrJSONMarshal:         ErrCodeJSONMarshal,
	ErrJSONUnmarshal:       ErrCodeJSONUnmarshal,
	ErrLongPaginationSize:  ErrCodeLongPaginationSize,
	ErrShortPaginationSize: ErrCodeShortPaginationSize,

	// Users
	ErrUserCount:  ErrCodeUserCount,
	ErrUserCreate: ErrCodeUserCreate,
	ErrUserDelete: ErrCodeUserDelete,
	ErrUserRead:   ErrCodeUserRead,
	ErrUserReads:  ErrCodeUserReads,
	ErrUserSearch: ErrCodeUserSearch,
	ErrUserUpdate: ErrCodeUserUpdate,

	// Posts
	ErrPostCount:    ErrCodePostCount,
	ErrPostCreate:   ErrCodePostCreate,
	ErrPostDelete:   ErrCodePostDelete,
	ErrPostRead:     ErrCodePostRead,
	ErrPostReads:    ErrCodePostReads,
	ErrPostUpdate:   ErrCodePostUpdate,
	ErrPostNotFound: ErrCodePostNotFound,

	// Others
	ErrFailedRead:        ErrCodeFailedRead,
	ErrFailedSave:        ErrCodeFailedSave,
	ErrInvalidPassword:   ErrCodeInvalidPassword,
	ErrInvalidQueryParam: ErrCodeInvalidQueryParam,
	ErrInvalidRequest:    ErrCodeInvalidRequest,
	ErrInvalidRole:       ErrCodeInvalidRole,
	ErrUnauthorized:      ErrCodeUnauthorized,
}

// Code gets machine-readable error code from error.
func Code(err error) string {
	if code, ok := errorCodes[err]; ok {
		return code
	}

	return err.Error()
}

var statusCodeMap = map[string]int{
	// Auth
	ErrCodeEmailAlreadyTaken:    http.StatusBadRequest,
	ErrCodeEmailRequired:        http.StatusBadRequest,
	ErrCodeInvalidToken:         http.StatusUnauthorized,
	ErrCodeExpiredToken:         http.StatusUnauthorized,
	ErrCodeInvalidRequest:       http.StatusBadRequest,
	ErrCodeLongPassword:         http.StatusBadRequest,
	ErrCodeLongUsername:         http.StatusBadRequest,
	ErrCodeLoginFailed:          http.StatusBadRequest,
	ErrCodeMissingAuthHeader:    http.StatusUnauthorized,
	ErrCodeShortPassword:        http.StatusBadRequest,
	ErrCodeShortUsername:        http.StatusBadRequest,
	ErrCodeUnauthorized:         http.StatusUnauthorized,
	ErrCodeUserDisabled:         http.StatusUnauthorized,
	ErrCodeUserNotFound:         http.StatusNotFound,
	ErrCodeEmailNotFound:        http.StatusBadRequest,
	ErrCodeUsernameAlreadyTaken: http.StatusBadRequest,
	ErrCodeUsernameRequired:     http.StatusBadRequest,
	ErrCodeWeakPassword:         http.StatusBadRequest,

	// Common
	ErrCodeBadRequest:           http.StatusBadRequest,
	ErrCodeBinding:              http.StatusBadRequest,
	ErrCodeCollectionIDRequired: http.StatusBadRequest,
	ErrCodeDocumentNotFound:     http.StatusNotFound,
	ErrCodeEmptyID:              http.StatusBadRequest,
	ErrCodeInvalidID:            http.StatusBadRequest,
	ErrCodeJSONDecode:           http.StatusUnprocessableEntity,
	ErrCodeJSONEncode:           http.StatusUnprocessableEntity,
	ErrCodeJSONMarshal:          http.StatusUnprocessableEntity,
	ErrCodeJSONUnmarshal:        http.StatusUnprocessableEntity,
	ErrCodeLongCollectionID:     http.StatusBadRequest,
	ErrCodeLongPaginationSize:   http.StatusBadRequest,
	ErrCodeShortCollectionID:    http.StatusBadRequest,
	ErrCodeShortPaginationSize:  http.StatusBadRequest,
	ErrCodeURLInvalid:           http.StatusBadRequest,
	ErrCodeURLRequired:          http.StatusBadRequest,
	ErrCodeUserAgentReadFile:    http.StatusUnprocessableEntity,

	// User
	ErrCodeUserCount:  http.StatusUnprocessableEntity,
	ErrCodeUserCreate: http.StatusUnprocessableEntity,
	ErrCodeUserDelete: http.StatusUnprocessableEntity,
	ErrCodeUserRead:   http.StatusUnprocessableEntity,
	ErrCodeUserReads:  http.StatusUnprocessableEntity,
	ErrCodeUserSearch: http.StatusUnprocessableEntity,
	ErrCodeUserUpdate: http.StatusUnprocessableEntity,

	// Post
	ErrCodePostCount:    http.StatusUnprocessableEntity,
	ErrCodePostCreate:   http.StatusUnprocessableEntity,
	ErrCodePostDelete:   http.StatusUnprocessableEntity,
	ErrCodePostRead:     http.StatusUnprocessableEntity,
	ErrCodePostReads:    http.StatusUnprocessableEntity,
	ErrCodePostUpdate:   http.StatusUnprocessableEntity,
	ErrCodePostNotFound: http.StatusNotFound,
}

// StatusCode gets HTTP status code from error code.
func StatusCode(code string) int {
	if status, ok := statusCodeMap[code]; ok {
		return status
	}

	return http.StatusInternalServerError
}

func Required(field string) error {
	switch field {
	case "Email":
		return ErrEmailRequired
	case "Username":
		return ErrUsernameRequired
	}

	return fmt.Errorf("%s is required", field)
}

func Max(field string) error {
	switch field {
	case "Password":
		return ErrLongPassword
	case "Username":
		return ErrLongUsername
	case "Size":
		return ErrLongPaginationSize
	}

	return fmt.Errorf("%s too long", field)
}

func Min(field string) error {
	switch field {
	case "Password":
		return ErrShortPassword
	case "Username":
		return ErrShortUsername
	case "Size":
		return ErrShortPaginationSize
	}

	return fmt.Errorf("%s too short", field)
}

func Len(field string) error {
	switch field {
	case "ID":
		return ErrInvalidID
	}

	return fmt.Errorf("%s invalid", field)
}

func ValidationError(errors []*APIError) error {
	return &APIErrors{
		errors,
	}
}
