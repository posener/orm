package load

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	personTypeString = "struct{Name string; Age int; unexported bool}"
	pkgName          = "example"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	tests := []struct {
		pkg, name  string
		wantStruct string
		wantPkg    string
		wantErr    bool
	}{
		{
			pkg:        "github.com/posener/orm/example",
			name:       "Person",
			wantStruct: personTypeString,
			wantPkg:    pkgName,
		},
		{
			pkg:        "../example",
			name:       "Person",
			wantStruct: personTypeString,
			wantPkg:    pkgName,
		},
		{
			pkg:     "github.com/posener/orm/example",
			name:    "NoOne",
			wantErr: true,
		},
		{
			pkg:     "nosuchpkg",
			name:    "Person",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s.%s", tt.pkg, tt.name), func(t *testing.T) {
			st, err := Load(tt.pkg, tt.name)
			if assert.Equal(t, tt.wantErr, err != nil) && !tt.wantErr {
				assert.Equal(t, tt.wantStruct, st.Struct.String())
				assert.Equal(t, tt.wantPkg, st.Pkg.Name())
			}
		})
	}
}
