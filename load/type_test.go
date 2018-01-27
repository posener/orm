package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestType(t *testing.T) {
	tests := []struct {
		typeName       string
		wantString     string
		wantType       string
		wantImportPath string
		wantExtName    string
		wantExtNaked   string
		wantPackage    string
		wantPointer    bool
		wantSlice      bool
		wantIsBasic    bool
	}{
		{
			typeName:     "int64",
			wantType:     "int64",
			wantString:   "int64",
			wantExtName:  "int64",
			wantExtNaked: "int64",
			wantIsBasic:  true,
		},
		{
			typeName:     "*int64",
			wantType:     "*int64",
			wantString:   "*int64",
			wantExtName:  "*int64",
			wantExtNaked: "int64",
			wantPointer:  true,
			wantIsBasic:  true,
		},
		{
			typeName:     "[]int64",
			wantType:     "[]int64",
			wantString:   "[]int64",
			wantExtName:  "[]int64",
			wantExtNaked: "int64",
			wantSlice:    true,
			wantIsBasic:  true,
		},
		{
			typeName:     "[]*int64",
			wantType:     "[]*int64",
			wantString:   "[]*int64",
			wantExtName:  "[]*int64",
			wantExtNaked: "int64",
			wantPointer:  true,
			wantSlice:    true,
			wantIsBasic:  true,
		},
		{
			typeName:       "github.com/posener/orm/tests.Person",
			wantType:       "Person",
			wantImportPath: "github.com/posener/orm/tests",
			wantString:     "github.com/posener/orm/tests.Person",
			wantExtName:    "tests.Person",
			wantExtNaked:   "tests.Person",
			wantPackage:    "tests",
		},
		{
			typeName:       "*github.com/posener/orm/tests.Person",
			wantType:       "*Person",
			wantImportPath: "github.com/posener/orm/tests",
			wantString:     "*github.com/posener/orm/tests.Person",
			wantExtName:    "*tests.Person",
			wantExtNaked:   "tests.Person",
			wantPackage:    "tests",
			wantPointer:    true,
		},
		{
			typeName:       "[]github.com/posener/orm/tests.Person",
			wantType:       "[]Person",
			wantImportPath: "github.com/posener/orm/tests",
			wantString:     "[]github.com/posener/orm/tests.Person",
			wantExtName:    "[]tests.Person",
			wantExtNaked:   "tests.Person",
			wantPackage:    "tests",
			wantSlice:      true,
		},
		{
			typeName:       "[]*github.com/posener/orm/tests.Person",
			wantType:       "[]*Person",
			wantImportPath: "github.com/posener/orm/tests",
			wantString:     "[]*github.com/posener/orm/tests.Person",
			wantExtName:    "[]*tests.Person",
			wantExtNaked:   "tests.Person",
			wantPackage:    "tests",
			wantPointer:    true,
			wantSlice:      true,
		},
		{
			typeName:       "../tests.Person",
			wantType:       "Person",
			wantImportPath: "github.com/posener/orm/tests",
			wantString:     "github.com/posener/orm/tests.Person",
			wantExtName:    "tests.Person",
			wantExtNaked:   "tests.Person",
			wantPackage:    "tests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			tp, err := New(tt.typeName)
			require.Nil(t, err)
			assert.Equal(t, tt.wantString, tp.String())
			assert.Equal(t, tt.wantExtName, tp.Ext(""))
			assert.Equal(t, tt.wantPackage, tp.Package())
			assert.Equal(t, tt.wantImportPath, tp.ImportPath())
			assert.Equal(t, tt.wantExtNaked, tp.Type.Ext(""))
			assert.Equal(t, tt.wantPointer, tp.Pointer)
			assert.Equal(t, tt.wantSlice, tp.Slice)
			assert.Equal(t, tt.wantIsBasic, tp.IsBasic())
		})
	}
}
