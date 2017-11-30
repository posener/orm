package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		typeName        string
		wantString      string
		wantType        string
		wantImportPath  string
		wantExtTypeName string
		wantNonPointer  string
		wantPackage     string
		wantIsPointer   bool
		wantIsBasic     bool
	}{
		{
			typeName:        "int64",
			wantType:        "int64",
			wantString:      "int64",
			wantExtTypeName: "int64",
			wantNonPointer:  "int64",
			wantIsBasic:     true,
		},
		{
			typeName:        "*int64",
			wantType:        "*int64",
			wantString:      "*int64",
			wantExtTypeName: "*int64",
			wantNonPointer:  "int64",
			wantIsPointer:   true,
			wantIsBasic:     true,
		},
		{
			typeName:        "github.com/posener/orm/example.Person",
			wantType:        "Person",
			wantImportPath:  "github.com/posener/orm/example",
			wantString:      "github.com/posener/orm/example.Person",
			wantExtTypeName: "example.Person",
			wantNonPointer:  "example.Person",
			wantPackage:     "example",
		},
		{
			typeName:        "*github.com/posener/orm/example.Person",
			wantType:        "*Person",
			wantImportPath:  "github.com/posener/orm/example",
			wantString:      "*github.com/posener/orm/example.Person",
			wantExtTypeName: "*example.Person",
			wantNonPointer:  "example.Person",
			wantPackage:     "example",
			wantIsPointer:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			tp, err := New(tt.typeName)
			require.Nil(t, err)
			assert.Equal(t, tt.wantString, tp.String())
			assert.Equal(t, tt.wantExtTypeName, tp.ExtTypeName())
			assert.Equal(t, tt.wantPackage, tp.Package())
			assert.Equal(t, tt.wantImportPath, tp.ImportPath)
			assert.Equal(t, tt.wantNonPointer, tp.NonPointer())
			assert.Equal(t, tt.wantIsPointer, tp.IsPointer())
			assert.Equal(t, tt.wantIsBasic, tp.IsBasic())
		})
	}
}
