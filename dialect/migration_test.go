package dialect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable_Diff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		got      *Table
		want     *Table
		wantDiff *Table
		wantErr  bool
	}{
		{
			name: "simple",
			got:  &Table{Columns: []Column{{Name: "a"}}},
			want: &Table{
				Columns:     []Column{{Name: "a"}, {Name: "b"}, {Name: "c"}},
				PrimaryKeys: []string{"a"},
				ForeignKeys: []ForeignKey{{Columns: []string{"b", "c"}}},
			},
			wantDiff: &Table{
				Columns:     []Column{{Name: "b"}, {Name: "c"}},
				PrimaryKeys: []string{"a"},
				ForeignKeys: []ForeignKey{{Columns: []string{"b", "c"}}},
			},
		},
		{
			name: "same",
			got: &Table{
				Columns:     []Column{{Name: "a"}, {Name: "b"}, {Name: "c"}},
				PrimaryKeys: []string{"a"},
				ForeignKeys: []ForeignKey{{Columns: []string{"b", "c"}}},
			},
			want: &Table{
				Columns:     []Column{{Name: "a"}, {Name: "b"}, {Name: "c"}},
				PrimaryKeys: []string{"a"},
				ForeignKeys: []ForeignKey{{Columns: []string{"b", "c"}}},
			},
			wantDiff: &Table{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDiff, err := tt.got.Diff(tt.want)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.wantDiff, gotDiff)
			}
		})
	}
}
