package memeng

import (
	"sync"

	wwr "github.com/qbeon/webwire-go"
	eng "github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/engine"
)

// engine represents an in-memory implementation of the API Engine interface
type engine struct {
	// lock protects the engine internals from concurrent access
	// and prevents data races
	lock *sync.Mutex

	// sessions stores all currently open sessions indexed by key
	sessions map[string]*wwr.Session
}

// New initializes a new in-memory engine implementation
func New() eng.Engine {
	return &engine{
		lock:     &sync.Mutex{},
		sessions: make(map[string]*wwr.Session),
	}
}
