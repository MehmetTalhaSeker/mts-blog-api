package appcontext_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/MehmetTalhaSeker/mts-blog-api/internal/appcontext"
)

func TestWithCorrelationID(t *testing.T) {
	id := uuid.NewString()

	ctx := appcontext.WithCorrelationID(context.Background(), id)

	val, err := appcontext.CorrelationID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, id, val)
}

func TestCorrelationID(t *testing.T) {
	id := uuid.NewString()

	cases := map[string]struct {
		ctx     context.Context
		want    string
		wantErr bool
	}{
		"request ID missing": {
			ctx:     context.Background(),
			want:    "",
			wantErr: true,
		},
		"appcontext with request ID": {
			ctx:     appcontext.WithCorrelationID(context.Background(), id),
			want:    id,
			wantErr: false,
		},
	}

	for desc, tc := range cases {
		t.Run(desc, func(t *testing.T) {
			got, err := appcontext.CorrelationID(tc.ctx)
			assert.Equal(t, tc.wantErr, err != nil)
			assert.Equal(t, tc.want, got)
		})
	}
}
