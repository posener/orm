package load

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	tests := []struct {
		pkg, name      string
		wantName       string
		wantFields     []Field
		wantFullName   string
		wantImportPath string
		wantErr        bool
	}{
		{
			pkg:          "github.com/posener/orm/example",
			name:         "Person",
			wantName:     "Person",
			wantFullName: "example.Person",
			wantFields: []Field{
				{Name: "Name", GoType: GoType{Type: "string"}, SQL: SQL{Column: "name"}},
				{Name: "Age", GoType: GoType{Type: "int"}, SQL: SQL{Column: "age"}},
			},
			wantImportPath: "github.com/posener/orm/example",
		},
		{
			pkg:          "../example",
			name:         "Person",
			wantName:     "Person",
			wantFullName: "example.Person",
			wantFields: []Field{
				{Name: "Name", GoType: GoType{Type: "string"}, SQL: SQL{Column: "name"}},
				{Name: "Age", GoType: GoType{Type: "int"}, SQL: SQL{Column: "age"}},
			},
			wantImportPath: "github.com/posener/orm/example",
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
		{
			pkg:          "../example",
			name:         "Employee",
			wantName:     "Employee",
			wantFullName: "example.Employee",
			wantFields: []Field{
				{Name: "Name", GoType: GoType{Type: "string"}, SQL: SQL{Column: "name"}},
				{Name: "Age", GoType: GoType{Type: "int"}, SQL: SQL{Column: "age"}},
				{Name: "Salary", GoType: GoType{Type: "int"}, SQL: SQL{Column: "salary"}},
			},
			wantImportPath: "github.com/posener/orm/example",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s.%s", tt.pkg, tt.name), func(t *testing.T) {
			tp, err := Load(GoType{ImportPath: tt.pkg, Type: tt.name})
			if assert.Equal(t, tt.wantErr, err != nil) && !tt.wantErr {
				assert.Equal(t, tt.wantName, tp.Name())
				assert.Equal(t, tt.wantFullName, tp.ExtTypeName())
				assert.Equal(t, tt.wantFields, tp.Fields)
				assert.Equal(t, tt.wantImportPath, tp.ImportPath)
			}
		})
	}
}
