package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/romshark/gapi/compiler"
	"github.com/romshark/gapi/compiler/parser"
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

	// Compiler
	ast, err := compiler.Compile(parser.SourceFile{
		File: parser.File{
			Name: filepath.Base(*schemaFilePath),
			Path: filepath.Dir(*schemaFilePath),
		},
		Src: string(fileContents),
	})
	if err != nil {
		log.Fatalf("compiler: %s", err)
	}

	log.Print("SUCCESS: ", ast)
	log.Print("SCHEMA NAME: ", ast.SchemaName)
}
