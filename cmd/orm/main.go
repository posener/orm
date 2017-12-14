package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/posener/orm/gen"
	"github.com/posener/orm/graph"
	"github.com/posener/orm/load"
)

var (
	options struct {
		types stringSlice
	}
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(val string) error {
	*s = strings.Split(val, ",")
	return nil
}

func init() {
	flag.Var(&options.types, "type", `Type name.
        Might be of the form '<pkg>.<type>' or just '<type>'. Where:
        - <pkg> can be either package name (e.x 'github.com/repository/project/package')
          or relative path (e.x './' or '../package').
        - <type> should be the type name.`)
	flag.Parse()
}

func main() {
	if len(options.types) == 0 {
		log.Fatal("Must give type full name")
	}

	var errors []string

	for _, typeName := range options.types {
		log.Printf("Loading type")
		tp, err := load.New(typeName)
		if err != nil {
			errors = append(errors, fmt.Sprintf("[%s] load type: %s", typeName, err))
			continue
		}

		log.Printf("Calculating graph")
		g, err := graph.New(tp)
		if err != nil {
			errors = append(errors, fmt.Sprintf("[%s] setting relations: %s", typeName, err))
			continue
		}

		log.Printf("Generating code")
		err = gen.Gen(g)
		if err != nil {
			errors = append(errors, fmt.Sprintf("[%s] generate code: %s", typeName, err))
		}
	}
	if len(errors) != 0 {
		log.Fatalf("Failed:\n%s", strings.Join(errors, "\n"))
	} else {
		log.Printf("Finished successfully!")
	}
}
