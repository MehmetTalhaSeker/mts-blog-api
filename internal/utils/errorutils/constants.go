package errorutils

// Auth Error Codes.
const (
	ErrCodeEmailAlreadyTaken    = "auth/email-taken"
	ErrCodeEmailRequired        = "auth/email-required"
	ErrCodeInvalidToken         = "auth/invalid-token"
	ErrCodeLongPassword         = "auth/long-password"
	ErrCodeLongUsername         = "auth/long-username"
	ErrCodeMissingAuthHeader    = "auth/missing-header"
	ErrCodeShortPassword        = "auth/short-password"
	ErrCodeShortUsername        = "auth/short-username"
	ErrCodeUserDisabled         = "auth/user-disabled"
	ErrCodeUserNotFound         = "auth/user-not-found"
	ErrCodeUsernameAlreadyTaken = "auth/username-taken"
	ErrCodeUsernameRequired     = "auth/username-required"
	ErrCodeWeakPassword         = "auth/weak-password"
)

// Common Error Codes.
const (
	ErrCodeBadRequest           = "req/bad-request"
	ErrCodeBinding              = "req/binding"
	ErrCodeCollectionIDRequired = "com/collection-id-required"
	ErrCodeDocumentNotFound     = "com/doc-not-found"
	ErrCodeEmptyID              = "req/empty-id"
	ErrCodeInvalidID            = "com/invalid-id"
	ErrCodeJSONDecode           = "com/json-decode"
	ErrCodeJSONEncode           = "com/json-encode"
	ErrCodeJSONMarshal          = "com/json-marshal"
	ErrCodeJSONUnmarshal        = "com/json-unmarshal"
	ErrCodeLongCollectionID     = "com/long-collection-id-size"
	ErrCodeLongPaginationSize   = "com/long-pagination-size"
	ErrCodeShortCollectionID    = "com/short-collection-id-size"
	ErrCodeShortPaginationSize  = "com/short-pagination-size"
	ErrCodeURLInvalid           = "com/url-invalid"
	ErrCodeURLRequired          = "com/url-required"
	ErrCodeUserAgentReadFile    = "com/user-agent-read"
)

// User Error Codes.
const (
	ErrCodeUserCount  = "user/count-failed"
	ErrCodeUserCreate = "user/create-failed"
	ErrCodeUserDelete = "user/delete-failed"
	ErrCodeUserRead   = "user/read-failed"
	ErrCodeUserReads  = "user/reads-failed"
	ErrCodeUserSearch = "user/search-failed"
	ErrCodeUserUpdate = "user/update-failed"
)

// Unorganized Error Codes.
const (
	ErrCodeFailedRead              = "un/read-failed"
	ErrCodeFailedSave              = "un/save-failed"
	ErrCodeInvalidPasswordUsername = "un/invalid"
	ErrCodeInvalidQueryParam       = "un/"
	ErrCodeInvalidRequest          = "un/invalid-request"
	ErrCodeInvalidRole             = "un/"
	ErrCodeUnauthorized            = "un/"
)
