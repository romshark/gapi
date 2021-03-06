package parser

import (
	"fmt"
	"sync"
)

// File represents a source file
type File struct {
	Name string
	Path string
}

// SourceFile represents an input source file
type SourceFile struct {
	File
	Src string
}

// Parser represents a GAPI parser
type Parser struct {
	errors            []Error
	errorsLock        *sync.Mutex
	deferredJobs      []func()
	mod               *SchemaModel
	lastIssuedGraphID GraphNodeID
	lastIssuedTypeID  TypeID
	lastIssuedParamID ParamID
	graphNodeByName   map[string]GraphNode
	typeByName        map[string]Type
	typeByID          map[TypeID]Type
	graphNodeByID     map[GraphNodeID]GraphNode
	paramByID         map[ParamID]*Parameter
}

// NewParser creates a new GAPI parser instance
func NewParser() (*Parser, error) {
	return &Parser{
		errorsLock: &sync.Mutex{},
	}, nil
}

// ResetState resets the parser state
func (pr *Parser) ResetState() {
	pr.mod = nil
	pr.lastIssuedGraphID = 0
	pr.lastIssuedTypeID = TypeIDUserTypeOffset
	pr.lastIssuedParamID = 0
	pr.graphNodeByName = make(map[string]GraphNode)
	pr.typeByName = make(map[string]Type)
	pr.typeByID = make(map[TypeID]Type)
	pr.graphNodeByID = make(map[GraphNodeID]GraphNode)
	pr.paramByID = make(map[ParamID]*Parameter)
}

// deferJob defers a function up until the parser has finished scanning
func (pr *Parser) deferJob(job func()) {
	pr.deferredJobs = append(pr.deferredJobs, job)
}

// err logs a parser error returning true if an error was logged,
// otherwise returning false
func (pr *Parser) err(err Error) bool {
	if err == nil {
		return false
	}
	if err.Code() == 0 {
		panic("invalid error code (0)")
	}
	pr.errorsLock.Lock()
	pr.errors = append(pr.errors, err)
	pr.errorsLock.Unlock()
	return true
}

// Errors returns a copy of the list of all compiler errors
func (pr *Parser) Errors() []Error {
	errs := make([]Error, len(pr.errors))
	copy(errs, pr.errors)
	return errs
}

// SchemaModel returns a copy of the schema model or nil if parsing failed
// or wasn't yet executed
func (pr *Parser) SchemaModel() *SchemaModel {
	if len(pr.errors) > 0 {
		return nil
	}
	return pr.mod.Clone()
}

// Parse starts parsing the source code reseting the parser
func (pr *Parser) Parse(source SourceFile) error {
	pr.ResetState()
	wg := &sync.WaitGroup{}

	// Initialize the model
	pr.mod = &SchemaModel{
		Types:          make([]Type, 0),
		EnumTypes:      make([]Type, 0),
		UnionTypes:     make([]Type, 0),
		QueryEndpoints: make([]*Query, 0),
		Mutations:      make([]*Mutation, 0),
	}

	// Initialize the lexer
	lexer := NewLexer(source)

	// Parse file
	fileFrag := pr.parseScmFile(lexer)
	if fileFrag == nil {
		goto END
	}

	// Execute all deferred jobs
	for j := 0; j < len(pr.deferredJobs); j++ {
		pr.deferredJobs[j]()
	}

	// Sort everything by name (ascending)
	wg.Add(7)
	go func() { sortTypesByName(pr.mod.Types); wg.Done() }()
	go func() { sortTypesByName(pr.mod.EnumTypes); wg.Done() }()
	go func() { sortTypesByName(pr.mod.UnionTypes); wg.Done() }()
	go func() { sortTypesByName(pr.mod.StructTypes); wg.Done() }()
	go func() { sortTypesByName(pr.mod.ResolverTypes); wg.Done() }()
	go func() { sortQueryEndpointsByName(pr.mod.QueryEndpoints); wg.Done() }()
	go func() { sortMutationsByName(pr.mod.Mutations); wg.Done() }()
	//TODO: sort trait types
	wg.Wait()

	// Perform semantic analysis
	wg.Add(2)
	go func() {
		// Find all recursive alias type cycles
		defer wg.Done()
		cycles := pr.findAliasTypeCycles()
		for _, cycle := range cycles {
			pr.err(&pErr{
				code: ErrAliasRecurs,
				message: fmt.Sprintf(
					"Recursive alias type cycle: %s",
					cycle.String(),
				),
			})
		}
	}()
	go func() {
		// Find all recursive struct type cycles
		defer wg.Done()
		cycles := pr.findStructTypeCycles()
		for _, cycle := range cycles {
			pr.err(&pErr{
				code: ErrStructRecurs,
				message: fmt.Sprintf(
					"Recursive struct type cycle: %s",
					cycle.String(),
				),
			})
		}
	}()
	if len(pr.mod.QueryEndpoints) < 1 && len(pr.mod.Mutations) < 1 {
		pr.err(&pErr{
			code:    ErrNoEndpoints,
			message: fmt.Sprintf("The schema is missing API endpoints"),
		})
	}
	wg.Wait()

END:
	if len(pr.errors) > 0 {
		return ParseErr{pr.Errors()}
	}

	return nil
}
