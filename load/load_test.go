package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	tests := []struct {
		typeName       string
		wantName       string
		wantFields     []Field
		wantFullName   string
		wantImportPath string
		wantErr        bool
	}{
		{
			// IMPORTANT: this test must be first, otherwise the package is loaded and the
			// test will not return an error as it should
			typeName: "nosuchpkg.Person",
			wantErr:  true,
		},
		{
			typeName: "github.com/posener/orm/example.NoOne",
			wantErr:  true,
		},
		{
			typeName:     "github.com/posener/orm/example.Person",
			wantName:     "Person",
			wantFullName: "example.Person",
			wantFields: []Field{
				{VarName: "Name", Type: Type{Name: "string"}, SQL: SQL{Column: "name"}},
				{VarName: "Age", Type: Type{Name: "int"}, SQL: SQL{Column: "age"}},
			},
			wantImportPath: "github.com/posener/orm/example",
		},
		{
			typeName:     "../example.Person",
			wantName:     "Person",
			wantFullName: "example.Person",
			wantFields: []Field{
				{VarName: "Name", Type: Type{Name: "string"}, SQL: SQL{Column: "name"}},
				{VarName: "Age", Type: Type{Name: "int"}, SQL: SQL{Column: "age"}},
			},
			wantImportPath: "github.com/posener/orm/example",
		},
		{
			typeName:     "../example.Employee",
			wantName:     "Employee",
			wantFullName: "example.Employee",
			wantFields: []Field{
				{VarName: "Name", Type: Type{Name: "string"}, SQL: SQL{Column: "name"}},
				{VarName: "Age", Type: Type{Name: "int"}, SQL: SQL{Column: "age"}},
				{VarName: "Salary", Type: Type{Name: "int"}, SQL: SQL{Column: "salary"}},
			},
			wantImportPath: "github.com/posener/orm/example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			tp, err := New(tt.typeName)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				require.Nil(t, err)
				assert.Equal(t, tt.wantName, tp.Name)
				assert.Equal(t, tt.wantFullName, tp.ExtTypeName())
				assert.Equal(t, tt.wantFields, tp.Fields)
				assert.Equal(t, tt.wantImportPath, tp.ImportPath)
			}
		})
	}
}
