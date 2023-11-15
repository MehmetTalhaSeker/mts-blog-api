package rbac_test

import (
	"context"
	"testing"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/appcontext"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/dto"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/rbac"
	"github.com/MehmetTalhaSeker/mts-blog-api/internal/types"
)

func TestIsMe(t *testing.T) {
	r := rbac.New()

	testCases := []struct {
		name         string
		userID       uint64
		claimsUserID uint64
		isMe         bool
	}{
		{
			name:         "Same user IDs",
			userID:       123,
			claimsUserID: 123,
			isMe:         true,
		},
		{
			name:         "Different user IDs",
			userID:       123,
			claimsUserID: 456,
			isMe:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims := &dto.Claims{
				UID: tc.claimsUserID,
			}
			ctx := context.Background()
			ctx = appcontext.WithMtsBlogUser(ctx, claims)

			isMe := r.IsMe(ctx, tc.userID)
			if isMe != tc.isMe {
				t.Errorf("expected isMe to be %v, got %v", tc.isMe, isMe)
			}
		})
	}
}

func TestIsAuthorized(t *testing.T) {
	r := rbac.New()

	testCases := []struct {
		name         string
		userRole     types.Role
		isAuthorized bool
	}{
		{
			name:         "Admin is authorized",
			userRole:     types.Admin,
			isAuthorized: true,
		},
		{
			name:         "Registered is not authorized",
			userRole:     types.Registered,
			isAuthorized: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims := &dto.Claims{
				Role: tc.userRole,
			}
			ctx := context.Background()
			ctx = appcontext.WithMtsBlogUser(ctx, claims)

			isAuthorized := r.IsModAuthorized(ctx)
			if isAuthorized != tc.isAuthorized {
				t.Errorf("expected isAuthorized to be %v, got %v", tc.isAuthorized, isAuthorized)
			}
		})
	}
}
