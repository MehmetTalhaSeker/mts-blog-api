package appcontext

import (
	"context"
	"errors"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

var (
	ErrInvalidCorrelationID = errors.New("unable to extract correlation ID from appcontext")
	ErrInvalidRole          = errors.New("unable to extract role from appcontext")
	ErrInvalidUser          = errors.New("unable to extract user from appcontext")
	ErrInvalidUserID        = errors.New("unable to extract user ID from appcontext")
	ErrInvalidLanguage      = errors.New("unable to extract language from appcontext")
)

// mtsBlogUserCtxKey is the appcontext key for the user value.
type mtsBlogUserCtxKey struct{}

// mtsBlogRoleCtxKey is the appcontext key for the role value.
type mtsBlogRoleCtxKey struct{}

// mtsBlogUserIDCtxKey is the appcontext key for the user ID value.
type mtsBlogUserIDCtxKey struct{}

// langCtxKey is the appcontext key for the client language value.
type langCtxKey struct{}

// correlationIDCtxKey is the key to use when getting or setting appcontext ID in appcontext.
// A package-specific type is used to prevent conflict with keys used by other packages
// as per https://pkg.go.dev/context#WithValue.
type correlationIDCtxKey struct{}

// WithCorrelationID returns a copy of the appcontext with the given string set as the correlationID
// value, using the appcontext key - correlationIDCtxKey - specific to this package.
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDCtxKey{}, correlationID)
}

// CorrelationID returns the appcontext ID from the given appcontext (looked up
// using the appcontext key - correlationIDCtxKey - specific to this package) or an error.
// An error is returned if no value is associated with the key.
func CorrelationID(ctx context.Context) (string, error) {
	cid, ok := ctx.Value(correlationIDCtxKey{}).(string)
	if !ok {
		return "", ErrInvalidCorrelationID
	}

	return cid, nil
}

// WithMtsBlogUser associates the user value with the mtsBlogUserCtxKey in the given appcontext.
func WithMtsBlogUser(ctx context.Context, user *dto.Claims) context.Context {
	return context.WithValue(ctx, mtsBlogUserCtxKey{}, user)
}

// WithMtsBlogRole associates the role value with the mtsBlogRoleCtxKey in the given appcontext.
func WithMtsBlogRole(ctx context.Context, role types.Role) context.Context {
	return context.WithValue(ctx, mtsBlogRoleCtxKey{}, role)
}

// WithMtsBlogUserID associates the userID value with the mtsBlogUserIDCtxKey in the given appcontext.
func WithMtsBlogUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, mtsBlogUserIDCtxKey{}, userID)
}

// WithLang associates the language value with the langCtxKey in the given appcontext.
func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langCtxKey{}, lang)
}

// MtsBlogUser returns the user value associated with the mtsBlogUserCtxKey in the given appcontext.
func MtsBlogUser(ctx context.Context) (*dto.Claims, error) {
	user, ok := ctx.Value(mtsBlogUserCtxKey{}).(*dto.Claims)
	if !ok {
		return nil, ErrInvalidUser
	}

	return user, nil
}

// MtsBlogRole returns the role value associated with the mtsBlogRoleCtxKey in the given appcontext.
func MtsBlogRole(ctx context.Context) (types.Role, error) {
	role, ok := ctx.Value(mtsBlogRoleCtxKey{}).(types.Role)
	if !ok {
		return "", ErrInvalidRole
	}

	return role, nil
}

// MtsBlogUserID returns the userID value associated with the mtsBlogUserIDCtxKey in the given appcontext.
func MtsBlogUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(mtsBlogUserIDCtxKey{}).(string)
	if !ok {
		return "", ErrInvalidUserID
	}

	return userID, nil
}

// Lang returns the language value associated with the langCtxKey in the given appcontext.
func Lang(ctx context.Context) (string, error) {
	lang, ok := ctx.Value(langCtxKey{}).(string)
	if !ok {
		return "", ErrInvalidLanguage
	}

	return lang, nil
}
