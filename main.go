package main

import (
	"flag"
	"log"

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
}

func main() {
	if options.name == "" {
		log.Fatalf("Must give struct name")
	}
	pkg, tp, err := load.Load(options.pkg, options.name)
	if err != nil {
		log.Fatal(err)
	}
	println(pkg, tp)
}
