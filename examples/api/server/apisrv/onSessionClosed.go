package apisrv

import "github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/metrics"

// OnSessionClosed implements the wwr.SessionManager interface
func (srv *apiServer) OnSessionClosed(sessionKey string) error {
	// Perform engine call
	srv.engine.CloseSession(sessionKey)

	// Log session closure to metrics
	metrics.SessionClosed()

	return nil
}
