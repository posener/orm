package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		typeName       string
		wantName       string
		wantFields     []Field
		wantFullName   string
		wantImportPath string
		wantErr        bool
	}{
		{
			typeName: "github.com/posener/orm/example.NoOne",
			wantErr:  true,
		},
		{
			typeName:     "github.com/posener/orm/example.Person",
			wantName:     "Person",
			wantFullName: "example.Person",
			wantFields: []Field{
				{Name: "Name", Type: Type{Name: "string"}},
				{Name: "Age", Type: Type{Name: "int"}},
			},
			wantImportPath: "github.com/posener/orm/example",
		},
		{
			typeName:     "../example.Person",
			wantName:     "Person",
			wantFullName: "example.Person",
			wantFields: []Field{
				{Name: "Name", Type: Type{Name: "string"}},
				{Name: "Age", Type: Type{Name: "int"}},
			},
			wantImportPath: "github.com/posener/orm/example",
		},
		{
			typeName:     "../example.Employee",
			wantName:     "Employee",
			wantFullName: "example.Employee",
			wantFields: []Field{
				{Name: "Name", Type: Type{Name: "string"}},
				{Name: "Age", Type: Type{Name: "int"}},
				{Name: "Salary", Type: Type{Name: "int"}},
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
				assert.Equal(t, tt.wantFullName, tp.ExtName())
				assert.Equal(t, tt.wantFields, tp.Fields)
				assert.Equal(t, tt.wantImportPath, tp.ImportPath)
			}
		})
	}
}
