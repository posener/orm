package load

import (
	"errors"
	"go/types"
	"log"

	"go/importer"

	"golang.org/x/tools/go/loader"
)

// Struct is information about loaded struct
type Struct struct {
	// Name is name of struct
	Name string
	// Struct describes the actual struct
	Struct *types.Struct
	// Pkg describes the package of the struct
	Pkg *types.Package
	// PkgMap describes all loaded packages
	PkgMap map[string]*types.Package
}

var (
	ErrTypeNotFound = errors.New("type was not found")
)

// Load loads a struct with name 'structName' in a package 'pkg'
// It returns descriptors for the package of the struct and the struct itself.
func Load(pkg, structName string) (*Struct, error) {
	l := loader.Config{
		AllowErrors:         true,
		TypeCheckFuncBodies: func(_ string) bool { return false },
		TypeChecker: types.Config{
			DisableUnusedImportCheck: true,
			Importer:                 importer.Default(),
		},
	}
	l.Import(pkg)
	p, err := l.Load()
	if err != nil {
		return nil, err
	}

	s := new(Struct)
	s.Name = structName

	s.PkgMap = map[string]*types.Package{}
	for pkg, _ := range p.AllPackages {
		s.PkgMap[pkg.Name()] = pkg
	}

	for pkgName, pkg := range p.Imported {
		log.Printf("scanning package: %s", pkgName)
		for _, scope := range pkg.Scopes {
			st := lookup(scope.Parent(), structName)
			if st != nil {
				s.Struct = st
				s.Pkg = pkg.Pkg
				return s, nil
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
