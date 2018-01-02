package load

import (
	"errors"
	"fmt"
	"go/importer"
	"go/types"
	"sync"

	"golang.org/x/tools/go/loader"
)

var (
	// ErrTypeNotFound is returned when the request type was not found
	ErrTypeNotFound = errors.New("type was not found")

	loadConfig = loader.Config{
		AllowErrors:         true,
		TypeCheckFuncBodies: func(_ string) bool { return false },
		TypeChecker: types.Config{
			DisableUnusedImportCheck: true,
			Importer:                 importer.Default(),
		},
	}
	importCache = map[string]*loader.Program{}
	cacheLock   sync.Mutex
)

func Program(importPath string) (*loader.Program, error) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	if p := importCache[importPath]; p != nil {
		return p, nil
	}
	loadConfig.Import(importPath)
	p, err := loadConfig.Load()
	if err != nil {
		return nil, err
	}
	importCache[importPath] = p
	return p, err
}

// loadStruct loads struct information from go package
func (t *Naked) loadStruct(importPath string) error {
	structName := t.Name

	// if import path is not define, try to import from local directory
	if importPath == "" {
		importPath = "./"
	}
	p, err := Program(importPath)
	if err != nil {
		return fmt.Errorf("load program: %s", err)
	}

	for _, pkg := range p.Imported {
		for _, scope := range pkg.Scopes {
			st := lookup(scope.Parent(), structName)
			if st != nil {
				t.st = st
				t.pkg = pkg.Pkg
				return nil
			}
		}
	}
	return ErrTypeNotFound
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
