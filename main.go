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
	flag.StringVar(&options.pkg, "pkg", "./example", "package of struct")
	flag.StringVar(&options.name, "name", "Person", "struct name")
	flag.Parse()
}

func main() {
	if options.name == "" {
		log.Fatal("Must give struct name")
	}
	pkg, tp, err := load.Load(options.pkg, options.name)
	failOnErr(err, "load struct")

	failOnErr(gen.Gen(pkg, tp, options.name), "generating")
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatal(msg + ": " + err.Error())
}
