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
	ast               *AST
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
	pr.ast = nil
	pr.lastIssuedTypeID = TypeIDUserTypeOffset
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

// AST returns a copy of the abstract syntax tree or nil if parsing failed
// or wasn't yet executed
func (pr *Parser) AST() *AST {
	if len(pr.errors) > 0 {
		return nil
	}
	return pr.ast.Clone()
}

// Parse starts parsing the source code reseting the parser
func (pr *Parser) Parse(source SourceFile) error {
	pr.ResetState()
	wg := &sync.WaitGroup{}

	// Initialize AST
	pr.ast = &AST{
		Types:          make([]Type, 0),
		AliasTypes:     make([]Type, 0),
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
	for _, job := range pr.deferredJobs {
		job()
	}

	// Sort everything by name (ascending)
	wg.Add(8)
	go func() { sortTypesByName(pr.ast.Types); wg.Done() }()
	go func() { sortTypesByName(pr.ast.AliasTypes); wg.Done() }()
	go func() { sortTypesByName(pr.ast.EnumTypes); wg.Done() }()
	go func() { sortTypesByName(pr.ast.UnionTypes); wg.Done() }()
	go func() { sortTypesByName(pr.ast.StructTypes); wg.Done() }()
	go func() { sortTypesByName(pr.ast.ResolverTypes); wg.Done() }()
	go func() { sortQueryEndpointsByName(pr.ast.QueryEndpoints); wg.Done() }()
	go func() { sortMutationsByName(pr.ast.Mutations); wg.Done() }()
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
	wg.Wait()

END:
	if len(pr.errors) > 0 {
		return ParseErr{pr.Errors()}
	}

	return nil
}
