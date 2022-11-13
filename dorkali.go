package dorkali

import (
	"fmt"
	"net/http"
)

type Result interface {
	// Title returns title of result
	Title() string

	// Description returns description of result
	Description() string

	// Url returns url of result
	Url() string

	// String returns result as string
	String() string
}

type Engine interface {
	// Use(...) call it when want to use it
	Start() error

	// Version returns engine version
	Version() string

	// Description of engine
	Description() string

	// Usage prints usage of engine
	Usage()

	// SearchContext searchs query and returns response
	Search(query interface{}) (*http.Response, error)

	// ParseResponse parses returned response from .Search(...) method
	ParseResponse(response *http.Response) ([]Result, error)

	// ParseHTML parses html responsed from google
	ParseHTML(h string) ([]Result, error)
}

var engines = make(map[string]func() Engine)

// register engine
func RegisterEngine(name string, new_engine func() Engine) {
	switch name {
	case "version", "help", "list":
		panic(name + ` is a engine name? you dont use these names: "version", "help", "list"`)
	}

	engines[name] = new_engine
}

// Engines returns registered engines names
func Engines() []string {
	keys := make([]string, 0, len(engines))
	for k := range engines {
		keys = append(keys, k)
	}
	return keys
}

type API struct {
	name string
	e    Engine
}

func (a *API) Name() string {
	return a.name
}

func (a *API) Start() error {
	return a.e.Start()
}

// Version returns engine version
func (a *API) Version() string {
	return a.e.Version()
}

// Search searchs query and returns response
func (a *API) Search(query interface{}) (*http.Response, error) {
	return a.e.Search(query)
}

// ParseResponse parses returned response from .SearchContext(...) or .Search(...) methods
func (a *API) ParseResponse(response *http.Response) ([]Result, error) {
	return a.e.ParseResponse(response)
}

// ParseHTML parses html responsed from google
func (a *API) ParseHTML(h string) ([]Result, error) {
	return a.e.ParseHTML(h)
}

func (a *API) Usage() {
	a.e.Usage()
}

// String returns string ( format: "API( name version | description )" )
func (a API) String() string {
	return "API( " + a.name + " " + a.e.Version() + " | " + a.e.Description() + " )"
}

// Use(...) returns engine if registered
//
// returns error if not found
//
// returns error if engine.Start() returns error
func Use(engineName string) (*API, error) {
	f, exists := engines[engineName]
	if !exists {
		return nil, fmt.Errorf("dorkali: unknown engine %q (forgotten import?)", engineName)
	}

	e := f()

	if err := e.Start(); err != nil {
		return nil, err
	}

	return &API{engineName, e}, nil
}

// UseWithoutStart(...) like Use(...) but not call Start()
//
// returns error if not found
func UseWithoutStart(engineName string) (*API, error) {
	f, exists := engines[engineName]
	if !exists {
		return nil, fmt.Errorf("dorkali: unknown engine %q (forgotten import?)", engineName)
	}

	return &API{engineName, f()}, nil
}
