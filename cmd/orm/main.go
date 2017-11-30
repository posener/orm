package main

import (
	"flag"
	"log"

	"github.com/posener/orm/gen"
	"github.com/posener/orm/load"
)

var (
	options struct {
		typeName string
	}
)

func init() {
	flag.StringVar(&options.typeName, "type", "", `Type name.
	Might be of the form '<pkg>.<type>' or just '<type>'. Where:
	- <pkg> can be either package name (e.x 'github.com/repository/project/package')
	  or relative path (e.x './' or '../package').
	- <type> should be the type name.`)
	flag.Parse()
}

func main() {
	if options.typeName == "" {
		log.Fatal("Must give type full name")
	}
	tp, err := load.New(options.typeName)
	failOnErr(err, "load type")
	failOnErr(gen.Gen(tp), "generating")
}

func failOnErr(err error, msg string) {
	if err == nil {
		return
	}
	log.Fatal(msg + ": " + err.Error())
}
