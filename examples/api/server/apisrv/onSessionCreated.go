package apisrv

import (
	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/metrics"
)

// OnSessionCreated implements the wwr.SessionManager interface
func (srv *apiServer) OnSessionCreated(conn wwr.Connection) error {
	// Perform engine call
	err := srv.engine.SaveSession(conn.Session())
	if err != nil {
		return err
	}

	// Log session creation to metrics
	metrics.SessionCreated()

	return nil
}
