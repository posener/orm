package graph

import (
	"fmt"
	"log"

	"github.com/posener/orm/load"
)

// Graph describes a relations graph from a root type.
// This graph has edges only from and to the root node.
// It could have more than one edge, in any direction between the root
// node and any other node.
type Graph struct {
	// Type is the root type
	*load.Type
	// In are relations from other types to the root type
	// Out are relations from the root type to other types
	// RelTable are relations from a relation table to the root type
	In, Out, RelTable []Edge
}

// Edge describes an edge in the graph
type Edge struct {
	// Field is the field in the root type from/to which the edge goes
	Field *load.Field
	// SrcField is the source field from which the edge goes from.
	// It can be in any node.
	SrcField     *load.Field
	RelationType *load.Type
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
			edge := Edge{
				Field:        field,
				RelationType: &field.Type,
			}
			log.Printf("one to one relation: by field %s to type %s", edge.Field, edge.Field.Type.Ext(""))
			root.Out = append(root.Out, edge)

		case field.IsReversedReference():
			srcField, err := findReversedSrcField(field)
			if err != nil {
				return nil, fmt.Errorf("field %v: %v", field, err)
			}

			if srcField != nil {
				edge := Edge{
					Field:        field,
					SrcField:     srcField,
					RelationType: tp,
				}
				log.Printf("many to one relation: Field %s is pointed by field %s", edge.Field, edge.SrcField)
				root.In = append(root.In, edge)
			} else {
				edge := Edge{
					Field:        field,
					RelationType: tp,
				}
				log.Printf("many relation: Field %s", edge.Field)
				root.RelTable = append(root.RelTable, edge)
			}
		}
	}
	return root, nil
}

// findReversedSrcField finds a field in a type that is referenced by the given reversedField.
// This field can be a specific one if RelationField is defined,
// or a field that points to the root type.
// In case of ambiguity, when more than one field is pointing to the root type,
// and non of them was given as the RelationField, an error will be returned.
func findReversedSrcField(reverseField *load.Field) (*load.Field, error) {
	var srcField *load.Field
	if reverseField.CustomRelationName != "" {
		return nil, nil
	}
	for _, field := range reverseField.Type.References() {
		if field.Type.Naked.String() != reverseField.ParentType.String() {
			continue
		}
		if !field.IsForwardReference() {
			continue
		}
		// if relation field name is defined, choose only this specific field
		if reverseField.RelationField == field.Name() {
			return field, nil
		}
		if srcField != nil {
			return nil, fmt.Errorf(
				"found more than one relation from %s to %s, please specify 'relation field' tag",
				reverseField.Type.Naked, reverseField.ParentType)
		}
		srcField = field
	}
	return srcField, nil
}
