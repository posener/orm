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
		wantFields     []*Field
		wantFullName   string
		wantLocalName  string
		wantImportPath string
		wantErr        bool
	}{
		{
			typeName: "github.com/posener/orm/tests.NoOne",
			wantErr:  true,
		},
		{
			typeName:      "github.com/posener/orm/tests.Person",
			wantName:      "Person",
			wantFullName:  "tests.Person",
			wantLocalName: "Person",
			wantFields: []*Field{
				{AccessName: "Name", Type: Type{Naked: &Naked{Name: "string"}}},
				{AccessName: "Age", Type: Type{Naked: &Naked{Name: "int"}}},
			},
			wantImportPath: "github.com/posener/orm/tests",
		},
		{
			typeName:      "../tests.Person",
			wantName:      "Person",
			wantFullName:  "tests.Person",
			wantLocalName: "Person",
			wantFields: []*Field{
				{AccessName: "Name", Type: Type{Naked: &Naked{Name: "string"}}},
				{AccessName: "Age", Type: Type{Naked: &Naked{Name: "int"}}},
			},
			wantImportPath: "github.com/posener/orm/tests",
		},
		{
			typeName:      "../tests.Employee",
			wantName:      "Employee",
			wantFullName:  "tests.Employee",
			wantLocalName: "Employee",
			wantFields: []*Field{
				{AccessName: "Person.Name", Type: Type{Naked: &Naked{Name: "string"}}},
				{AccessName: "Person.Age", Type: Type{Naked: &Naked{Name: "int"}}},
				{AccessName: "Salary", Type: Type{Naked: &Naked{Name: "int"}}},
			},
			wantImportPath: "github.com/posener/orm/tests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.typeName, func(t *testing.T) {
			tp, err := New(tt.typeName)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				require.Nil(t, err)
				require.Nil(t, tp.LoadFields(1))
				assert.Equal(t, tt.wantName, tp.Name)
				assert.Equal(t, tt.wantFullName, tp.Ext(""))
				assert.Equal(t, tt.wantLocalName, tp.Ext(tp.Package()))
				for _, f := range tp.Fields {
					f.ParentType = nil
				}
				assert.Equal(t, tt.wantFields, tp.Fields)
				assert.Equal(t, tt.wantImportPath, tp.ImportPath())
			}
		})
	}
}
