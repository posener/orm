package main

import (
	"flag"
	"log"

	"github.com/posener/orm/gen"
	"github.com/posener/orm/load"
)

var (
	options struct {
		pkg     string
		name    string
		dialect string
	}
)

func init() {
	flag.StringVar(&options.pkg, "pkg", ".", "package of struct")
	flag.StringVar(&options.name, "name", "", "struct name")
	flag.StringVar(&options.dialect, "dialect", "", "dialect of ORM")
	flag.Parse()
}

func main() {
	if options.name == "" {
		log.Fatal("Must give struct name")
	}
	if options.dialect == "" {
		log.Fatal("Must give dialect")
	}
	st, err := load.Load(options.pkg, options.name)
	failOnErr(err, "load struct")
	failOnErr(gen.Gen(st, options.dialect), "generating")
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatal(msg + ": " + err.Error())
}
