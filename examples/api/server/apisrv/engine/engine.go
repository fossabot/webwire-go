package engine

import (
	wwr "github.com/qbeon/webwire-go"
)

// Engine defines the service engine interface
type Engine interface {
	SaveSession(newSession *wwr.Session) error
	FindSession(key string) (wwr.SessionLookupResult, error)
	CloseSession(key string) error
}
