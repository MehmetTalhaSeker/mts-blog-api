package errorutils

// Auth Error Codes.
const (
	ErrCodeEmailAlreadyTaken    = "auth/email-taken"
	ErrCodeEmailRequired        = "auth/email-required"
	ErrCodeInvalidToken         = "auth/invalid-token"
	ErrCodeExpiredToken         = "auth/expired-token"
	ErrCodeLongPassword         = "auth/long-password"
	ErrCodeLongUsername         = "auth/long-username"
	ErrCodeLoginFailed          = "auth/login-failed"
	ErrCodeMissingAuthHeader    = "auth/missing-header"
	ErrCodeShortPassword        = "auth/short-password"
	ErrCodeShortUsername        = "auth/short-username"
	ErrCodeUserDisabled         = "auth/user-disabled"
	ErrCodeEmailNotFound        = "auth/email-not-found"
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
	ErrCodeUserCount    = "user/count-failed"
	ErrCodeUserCreate   = "user/create-failed"
	ErrCodeUserDelete   = "user/delete-failed"
	ErrCodeUserRead     = "user/read-failed"
	ErrCodeUserReads    = "user/reads-failed"
	ErrCodeUserSearch   = "user/search-failed"
	ErrCodeUserUpdate   = "user/update-failed"
	ErrCodeUserNotFound = "user/user-not-found"
)

// Post Error Codes.
const (
	ErrCodePostCount    = "post/count-failed"
	ErrCodePostCreate   = "post/create-failed"
	ErrCodePostDelete   = "post/delete-failed"
	ErrCodePostRead     = "post/read-failed"
	ErrCodePostReads    = "post/reads-failed"
	ErrCodePostUpdate   = "post/update-failed"
	ErrCodePostNotFound = "post/not-found"
)

// Comment Error Codes.
const (
	ErrCodeCommentCount    = "comment/count-failed"
	ErrCodeCommentCreate   = "comment/create-failed"
	ErrCodeCommentDelete   = "comment/delete-failed"
	ErrCodeCommentRead     = "comment/read-failed"
	ErrCodeCommentReads    = "comment/reads-failed"
	ErrCodeCommentNotFound = "comment/not-found"
)

// Unorganized Error Codes.
const (
	ErrCodeFailedRead        = "un/read-failed"
	ErrCodeFailedSave        = "un/save-failed"
	ErrCodeInvalidPassword   = "un/invalid"
	ErrCodeInvalidQueryParam = "un/"
	ErrCodeInvalidRequest    = "un/invalid-request"
	ErrCodeInvalidRole       = "un/"
	ErrCodeUnauthorized      = "un/"
)
