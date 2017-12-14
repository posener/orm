package graph

import (
	"fmt"

	"github.com/posener/orm/load"
)

// Graph describes a relations graph from a root type.
// This graph has edges only from and to the root node.
// It could have more than one edge, in any direction between the root
// node and any other node.
type Graph struct {
	// Type is the root type
	*load.Type
	// In and Out are edges to and from the root node
	In, Out []Edge
}

// Edge describes an edge in the graph
type Edge struct {
	// LocalField is the field in the root type from/to which the edge goes
	LocalField *load.Field
	// SrcField is the source field from which the edge goes from.
	// It can be in any node.
	SrcField *load.Field
}

// RelationType is the type that the edge relates to.
func (e *Edge) RelationType() *load.Type {
	return &e.LocalField.Type
}

// New returns a graph with a given root type
func New(tp *load.Type) (*Graph, error) {

	// Only 3 levels of field loaded from the root type are needed.
	// This makes sure that fields that are referencing the root type will have
	// their type loaded.
	// This does not include embedded types, which are loaded infinitely deep.
	err := tp.LoadFields(3)
	if err != nil {
		return nil, fmt.Errorf("loading fields: %s", err)
	}

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

// findReversedSrcField finds a field in a type that is referenced by the given reversedField.
// This field can be a specific one if ReferencingFieldName is defined,
// or a field that points to the root type.
// In case of ambiguity, when more than one field is pointing to the root type,
// and non of them was given as the ReferencingFieldName, an error will be returned.
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
