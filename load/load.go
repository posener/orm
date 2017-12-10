package load

import (
	"errors"
	"fmt"
	"go/importer"
	"go/types"
	"log"
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
	typeCache   = map[string]*Type{}
	cacheLock   sync.Mutex
)

func loadProgram(importPath string) (*loader.Program, error) {
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

// cacheGetOrUpdate get or updates the cache.
// the cache is used to prevent recursive load of fields
// - if the type exists in the cache, it return the full type and true value
// - if it does not exists, it sets it in the cache, return it and false value
// the bool return value means 'exists in cache'
func cacheGetOrUpdate(tp *Type) (*Type, bool) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	fullName := tp.String()
	if loaded := typeCache[fullName]; loaded != nil {
		return loaded, true
	}
	typeCache[fullName] = tp
	return tp, false
}

// loadStruct loads struct information from go package
func (t *Type) loadStruct() (*types.Struct, *types.Package, error) {
	log.Printf("Loading struct %s", t)
	structName := t.Name

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
