package load

import (
	"errors"
	"go/types"
	"log"

	"go/importer"

	"golang.org/x/tools/go/loader"
)

func Load(pkg, name string) (*types.Package, *types.Struct, error) {
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
		return nil, nil, err
	}
	for pkgName, pkg := range p.Imported {
		log.Printf("scanning package: %s", pkgName)
		for _, s := range pkg.Scopes {
			s := lookup(s.Parent(), name)
			if s != nil {
				return pkg.Pkg, s, nil
			}
		}
	}
	return nil, nil, errors.New("type was not found")
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
