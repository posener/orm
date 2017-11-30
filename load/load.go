package load

import (
	"errors"
	"go/importer"
	"go/types"
	"log"

	"sync"

	"fmt"

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
	loadMu sync.Mutex
)

func loadProgram(importPath string) (*loader.Program, error) {
	loadMu.Lock()
	defer loadMu.Unlock()
	loadConfig.Import(importPath)
	return loadConfig.Load()
}

// loadStruct loads struct information from go package
func (t *Type) loadStruct() (*types.Struct, *types.Package, error) {
	log.Printf("Loading struct %s", t)
	structName := t.nonPointerType()

	// if import path is not define, try to import from local directory
	importPath := t.ImportPath
	if importPath == "" {
		importPath = "./"
	}
	p, err := loadProgram(importPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loading program: %s", err)
	}

	for pkgName, pkg := range p.Imported {
		log.Printf("scanning package: %s", pkgName)
		for _, scope := range pkg.Scopes {
			st := lookup(scope.Parent(), structName)
			if st != nil {
				return st, pkg.Pkg, nil
			}
		}
	}
	return nil, nil, ErrTypeNotFound
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
