package main

import (
	"flag"
	"log"

	"github.com/posener/orm/gen"
	"github.com/posener/orm/load"
)

var (
	options struct {
		pkg  string
		name string
	}
)

func init() {
	flag.StringVar(&options.pkg, "pkg", ".", "package of struct")
	flag.StringVar(&options.name, "name", "", "struct name")
	flag.Parse()
}

func main() {
	if options.name == "" {
		log.Fatal("Must give struct name")
	}
	st, err := load.Load(load.GoType{Type: options.name, ImportPath: options.pkg})
	failOnErr(err, "load struct")
	failOnErr(gen.Gen(st), "generating")
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatal(msg + ": " + err.Error())
}
