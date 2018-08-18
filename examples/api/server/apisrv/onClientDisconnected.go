package apisrv

import (
	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/metrics"
)

// OnClientDisconnected implements the wwr.ServerImplementation interface
func (srv *apiServer) OnClientDisconnected(client wwr.Connection) {
	metrics.ConnectionClosed()
}
