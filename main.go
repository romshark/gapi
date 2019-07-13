package main

import (
	"flag"
	"io/ioutil"
	"log"
)

var schemaFilePath = flag.String("schema", "", "schema file path")

func main() {
	flag.Parse()

	log.Print("SCHEMA: ", *schemaFilePath)
	if *schemaFilePath == "" {
		log.Fatal("missing schema file path (use -schema)")
	}

	// load schema file
	fileContents, err := ioutil.ReadFile(*schemaFilePath)
	if err != nil {
		log.Fatalf("reading file: %s", err)
	}

	parser := GAPIParser{
		Buffer: string(fileContents),
	}

	if err := parser.Init(); err != nil {
		log.Fatalf("parser initialization: %s", err)
	}

	if err := parser.Parse(); err != nil {
		log.Fatalf("parser: %s", err)
	}

	log.Print("SUCCESS: ", parser.Tokens())
}
