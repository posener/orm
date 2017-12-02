package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		typeName       string
		wantString     string
		wantType       string
		wantImportPath string
		wantExtName    string
		wantNonPointer string
		wantPackage    string
		wantPointer    bool
		wantIsBasic    bool
	}{
		{
			typeName:       "int64",
			wantType:       "int64",
			wantString:     "int64",
			wantExtName:    "int64",
			wantNonPointer: "int64",
			wantIsBasic:    true,
		},
		{
			typeName:       "*int64",
			wantType:       "*int64",
			wantString:     "*int64",
			wantExtName:    "*int64",
			wantNonPointer: "int64",
			wantPointer:    true,
			wantIsBasic:    true,
		},
		{
			typeName:       "github.com/posener/orm/example.Person",
			wantType:       "Person",
			wantImportPath: "github.com/posener/orm/example",
			wantString:     "github.com/posener/orm/example.Person",
			wantExtName:    "example.Person",
			wantNonPointer: "example.Person",
			wantPackage:    "example",
		},
		{
			typeName:       "*github.com/posener/orm/example.Person",
			wantType:       "*Person",
			wantImportPath: "github.com/posener/orm/example",
			wantString:     "*github.com/posener/orm/example.Person",
			wantExtName:    "*example.Person",
			wantNonPointer: "example.Person",
			wantPackage:    "example",
			wantPointer:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			tp, err := New(tt.typeName)
			require.Nil(t, err)
			assert.Equal(t, tt.wantString, tp.String())
			assert.Equal(t, tt.wantExtName, tp.ExtName())
			assert.Equal(t, tt.wantPackage, tp.Package())
			assert.Equal(t, tt.wantImportPath, tp.ImportPath)
			assert.Equal(t, tt.wantNonPointer, tp.ExtNonPointer())
			assert.Equal(t, tt.wantPointer, tp.Pointer())
			assert.Equal(t, tt.wantIsBasic, tp.IsBasic())
		})
	}
}
