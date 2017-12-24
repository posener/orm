package gen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/posener/orm/dialect"
	"github.com/posener/orm/graph"
	"github.com/posener/orm/runtime/migration"
)

// TemplateData arguments for the templates
type TemplateData struct {
	// The name	of the new created package
	Graph    *graph.Graph // TODO: rename Type to Graph
	Dialects []dialect.API
	Public   string
	Private  string
	Table    *migration.Table
}

// Gen generates all the ORM files for a given struct in a given package.
// st is the type descriptor of the struct
func Gen(g *graph.Graph, dialects []dialect.API) error {
	// get the package ormDir on disk
	structPkgDir, err := packagePath(g.ImportPath())
	if err != nil {
		return err
	}

	data := TemplateData{
		Graph:    g,
		Dialects: dialects,
		Public:   g.Name,
		Private:  strings.ToLower(g.Name),
		Table:    migration.NewTable(g),
	}

	ormFileName := strings.ToLower(g.Name + "_orm.go")
	ormFilePath := filepath.Join(structPkgDir, ormFileName)

	log.Printf("Generating code for %s into %s", g.Type, ormFilePath)

	ormFile, err := os.Create(ormFilePath)
	if err != nil {
		return fmt.Errorf("creating file %s: %s", ormFilePath, err)
	}

	// write templates
	if err = tpl.Execute(ormFile, data); err != nil {
		return fmt.Errorf("executing template: %s", err)
	}
	format(ormFilePath)
	return nil
}

func packagePath(pkg string) (string, error) {
	for _, gopath := range filepath.SplitList(os.Getenv("GOPATH")) {
		pkgPath := filepath.Join(gopath, "src", pkg)
		f, err := os.Stat(pkgPath)
		if err == nil && f.IsDir() {
			return pkgPath, nil
		}
	}
	return "", fmt.Errorf("package path was not found: %s", pkg)
}

func format(path string) {
	_, err := exec.Command("goimports", "-w", path).CombinedOutput()
	if err != nil {
		log.Printf("Failed formatting package: %s", err)
	}
}
