package apisrv

import (
	"context"
	"time"

	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/metrics"
)

// OnRequest implements the wwr.ServerImplementation interface
func (srv *apiServer) OnRequest(
	_ context.Context,
	client wwr.Connection,
	message wwr.Message,
) (response wwr.Payload, err error) {
	startTime := time.Now()
	metrics.Request()

	// Handle request

	// Log request completion to metrics
	metrics.RequestCompleted(time.Since(startTime))

	// Reply to the request using the same data and encoding
	return message.Payload(), nil
}
