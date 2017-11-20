package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s    string
		want map[string]map[string]string
	}{
		{
			s: ``,
		},
		{
			s:    `key:"value"`,
			want: map[string]map[string]string{"key": {"value": ""}},
		},
		{
			s:    `key:"value:subvalue"`,
			want: map[string]map[string]string{"key": {"value": "subvalue"}},
		},
		{
			s:    `key:"value:subvalue;value2:subvalue2"`,
			want: map[string]map[string]string{"key": {"value": "subvalue", "value2": "subvalue2"}},
		},
		{
			s: `key:"value:subvalue;value2:subvalue2" key2 key3:"value3" key4:"v4;v5"`,
			want: map[string]map[string]string{
				"key":  {"value": "subvalue", "value2": "subvalue2"},
				"key2": {"": ""},
				"key3": {"value3": ""},
				"key4": {"v4": "", "v5": ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := Parse(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFields(t *testing.T) {
	t.Parallel()
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
			got := Fields(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
