package graph

import (
	"fmt"

	"github.com/posener/orm/load"
)

type Graph struct {
	*load.Type
	In, Out []Edge
}

func (g *Graph) Edges() []Edge {
	return append(g.In, g.Out...)
}

type Edge struct {
	SrcField   *load.Field
	LocalField *load.Field
}

func (e *Edge) RelationType() *load.Type {
	return &e.LocalField.Type
}

func New(tp *load.Type) (*Graph, error) {
	root := &Graph{Type: tp}

	for _, field := range tp.Fields {
		switch {
		case field.IsForwardReference():
			root.Out = append(root.Out, Edge{
				SrcField:   field,
				LocalField: field,
			})

		case field.IsReversedReference():
			srcField, err := findReversedSrcField(field)
			if err != nil {
				return nil, fmt.Errorf("field %v: %v", field, err)
			}
			root.In = append(root.In, Edge{
				SrcField:   srcField,
				LocalField: field,
			})
		}
	}
	return root, nil
}

func findReversedSrcField(reverseField *load.Field) (*load.Field, error) {
	var srcField *load.Field
	for _, field := range reverseField.Type.Fields {
		if !field.IsForwardReference() || field.Type.Naked.String() != reverseField.ParentType.String() {
			continue
		}
		// if referencing field name is defined, choose only this specific field
		if reverseField.ReferencingFieldName == field.Name() {
			return field, nil
		}
		if srcField != nil {
			return nil, fmt.Errorf(
				"found more than one relation from %s to %s, please specify 'referencing field' tag",
				reverseField.Type.Naked, reverseField.ParentType)
		}
		srcField = field
	}
	if srcField == nil {
		return nil, fmt.Errorf(
			"needed reversed reference from %v to %v was not found",
			reverseField.Type.Naked, reverseField.ParentType)
	}
	return srcField, nil
}
