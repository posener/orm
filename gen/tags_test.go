package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{
			s:    `a b c`,
			want: []string{"a", "b", "c"},
		},
		{
			s:    ` a bc   d  `,
			want: []string{"a", "bc", "d"},
		},
		{
			s:    ` a "b  c"  def`,
			want: []string{"a", `"b  c"`, "def"},
		},
		{
			s:    `aaa bbb"ccc ddd" fff"" ggg""hhh "" ""iii`,
			want: []string{`aaa`, `bbb"ccc ddd"`, `fff""`, `ggg""hhh`, `""`, `""iii`},
		},
		{
			s:    `aaa "`,
			want: []string{`aaa`, `"`},
		},
		{
			s:    `aaa "bbb ccc`,
			want: []string{`aaa`, `"bbb ccc`},
		},
		{
			s:    `aaa \"bbb ccc\"`,
			want: []string{`aaa`, `\"bbb`, `ccc\"`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := fields(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
