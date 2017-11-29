package load

import (
	"errors"
	"go/importer"
	"go/types"
	"log"

	"golang.org/x/tools/go/loader"
)

var (
	// ErrTypeNotFound is returned when the request type was not found
	ErrTypeNotFound = errors.New("type was not found")
)

// Load loads a struct with name 'structName' in a package 'pkg'
// It returns descriptors for the package of the struct and the struct itself.
func Load(goType GoType) (*Type, error) {
	log.Printf("Loading struct %s", goType)
	structName := goType.NonPointer()
	l := loader.Config{
		AllowErrors:         true,
		TypeCheckFuncBodies: func(_ string) bool { return false },
		TypeChecker: types.Config{
			DisableUnusedImportCheck: true,
			Importer:                 importer.Default(),
		},
	}
	l.Import(goType.ImportPath)
	p, err := l.Load()
	if err != nil {
		return nil, err
	}

	for pkgName, pkg := range p.Imported {
		log.Printf("scanning package: %s", pkgName)
		for _, scope := range pkg.Scopes {
			st := lookup(scope.Parent(), structName)
			if st != nil {
				tp := newType(structName, st, pkg.Pkg)
				return tp, nil
			}
		}
	}
	return nil, ErrTypeNotFound
}

// lookup is a recursive lookup for a struct with name in a scope
func lookup(s *types.Scope, name string) *types.Struct {
	if o := s.Lookup(name); o != nil {
		u := o.Type().Underlying()
		if s, ok := u.(*types.Struct); ok {
			return s
		}
	}
	for i := 0; i < s.NumChildren(); i++ {
		s := lookup(s.Child(i), name)
		if s != nil {
			return s
		}
	}
	return nil
}
