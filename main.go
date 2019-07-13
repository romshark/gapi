package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/romshark/gapi/compiler"
)

var schemaFilePath = flag.String("schema", "", "schema file path")

func main() {
	flag.Parse()

	log.Print("SCHEMA: ", *schemaFilePath)
	if *schemaFilePath == "" {
		log.Fatal("missing schema file path (use -schema)")
	}

	// Load schema file
	fileContents, err := ioutil.ReadFile(*schemaFilePath)
	if err != nil {
		log.Fatalf("reading file: %s", err)
	}

	// Initialize compiler
	compiler, err := compiler.NewCompiler()
	if err != nil {
		log.Fatalf("compiler init: %s", err)
	}

	// Compile
	ast, err := compiler.Compile(string(fileContents))
	if err != nil {
		log.Fatalf("compiler: %s", err)
	}

	log.Print("SUCCESS: ", ast)
	log.Print("SCHEMA NAME: ", ast.SchemaName)
}
