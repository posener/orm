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
	"github.com/posener/orm/load"
	"github.com/posener/orm/runtime/migration"
)

var goPaths []string

func init() {
	paths := filepath.SplitList(os.Getenv("GOPATH"))
	for _, path := range paths {
		goPaths = append(goPaths, filepath.Join(path, "src"))
	}
}

// TemplateData arguments for the templates
type TemplateData struct {
	// The name	of the new created package
	Graph    *graph.Graph
	Dialects []dialect.API
	Public   string
	Private  string
	Table    *migration.Table
	Package  string
}

// Gen generates all the ORM files for a given struct in a given package.
// st is the type descriptor of the struct
func Gen(g *graph.Graph, out string, dialects []dialect.API) error {
	// get the package ormDir on disk
	structPkgDir, err := packagePath(g.ImportPath())
	if err != nil {
		return fmt.Errorf("find struct package: %s", err)
	}
	outDir, err := outDir(out)
	if err != nil {
		return fmt.Errorf("create out directory: %s", err)
	}
	if outDir == "" {
		outDir = structPkgDir
	}

	outPkg, err := outPkg(outDir)
	if err != nil {
		return fmt.Errorf("find out package: %s", err)
	}

	data := TemplateData{
		Graph:    g,
		Dialects: dialects,
		Public:   g.Name,
		Private:  strings.ToLower(g.Name),
		Table:    migration.NewTable(g),
		Package:  outPkg,
	}

	ormFileName := strings.ToLower(g.Name + "_orm.go")
	ormFilePath := filepath.Join(outDir, ormFileName)

	log.Printf("Generating code for %s into %s", g.Type, ormFilePath)

	ormFile, err := os.Create(ormFilePath)
	if err != nil {
		return fmt.Errorf("create file %s: %s", ormFilePath, err)
	}

	// write templates
	if err = tpl.Execute(ormFile, data); err != nil {
		return fmt.Errorf("executing template: %s", err)
	}
	format(ormFilePath)
	return nil
}

func packagePath(pkg string) (string, error) {
	for _, gopath := range goPaths {
		pkgPath := filepath.Join(gopath, pkg)
		f, err := os.Stat(pkgPath)
		if err == nil && f.IsDir() {
			return pkgPath, nil
		}
	}
	return "", fmt.Errorf("package path was not found: %s", pkg)
}

func outDir(out string) (string, error) {
	if len(out) == 0 {
		return "", nil
	}
	var dir string
	switch {
	case out[0] == '.':
		// relative path
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(wd, out)
	default:
		// out is a package name
		var err error
		dir, err = packagePath(out)
		if err != nil {
			if len(goPaths) == 0 {
				return "", fmt.Errorf("for out with package path, must define GOPATH environemnt variable")
			}
			dir = filepath.Join(goPaths[0], out)
		}
	}
	// create the output directory
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return "", fmt.Errorf("creating out directory %s: %s", dir, err)
	}
	return dir, nil
}

func outPkg(outDir string) (string, error) {
	var pkgPath string
	for _, goPath := range goPaths {
		if strings.HasPrefix(outDir, goPath) {
			pkgPath = strings.Trim(outDir[len(goPath):], string(os.PathSeparator))
			break
		}
	}
	if pkgPath == "" {
		return "", fmt.Errorf("output direcotry not in any GOPATH")
	}

	outProgram, err := load.Program(pkgPath)
	if err != nil {
		return "", fmt.Errorf("load output package: %s", err)
	}
	outPkg := outProgram.Package(pkgPath)
	if outPkg == nil {
		return "", fmt.Errorf("find output package: %s", pkgPath)
	}

	pkgName := outPkg.Pkg.Name()
	if pkgName == "" {
		pkgName = filepath.Base(outDir)
	}

	return pkgName, nil
}

func format(path string) {
	_, err := exec.Command("goimports", "-w", path).CombinedOutput()
	if err != nil {
		log.Printf("Failed formatting package: %s", err)
	}
}
