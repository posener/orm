package sqltypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		s       string
		want    *Type
		wantErr bool
	}{
		{
			s:    "VARCHAR(42)",
			want: &Type{Name: "varchar", Size: 42},
		},
		{
			s:    "INTEGER",
			want: &Type{Name: "integer"},
		},
		{
			s:       "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got, err := New(tt.s)
			require.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
