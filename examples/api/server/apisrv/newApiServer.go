package apisrv

import (
	"log"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	wwr "github.com/qbeon/webwire-go"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/config"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/dam"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/engine"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/logger"
	"github.com/qbeon/webwire-go/examples/apiServer/server/apisrv/session"
)

// NewApiServer initializes a new API server instance
func NewApiServer(
	conf config.Config,
	engine engine.Engine,
) (ApiServer, error) {
	// Initialize loggers
	logger, err := logger.New(&conf)
	if err != nil {
		return nil, errors.Wrap(err, "logger initialization failed")
	}

	newApiServer := &apiServer{
		conf:   conf,
		lock:   &sync.RWMutex{},
		stop:   dam.New(1),
		log:    logger,
		engine: engine,
	}

	// Initialize webwire server
	newApiServer.wwrSrv, err = wwr.NewHeadlessServer(
		newApiServer,
		wwr.ServerOptions{
			// Enable webwire sessions
			Sessions: wwr.Enabled,

			// Make the API server responsible for handling the sessions
			SessionManager: newApiServer,

			// Define the session info parser
			SessionInfoParser: session.ParseSessionInfo,
			Heartbeat:         wwr.Enabled,

			// Use the log writers provided by the logger instance
			ErrorLog: log.New(
				logger.ErrorLogWriter(),
				"WWR_ERR: ",
				log.Ldate|log.Ltime|log.Lshortfile,
			),
			WarnLog: log.New(
				logger.ErrorLogWriter(),
				"WWR_WARN: ",
				log.Ldate|log.Ltime|log.Lshortfile,
			),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "webwire server initialization failed")
	}

	// Initialize the example HTTP endpoint
	newApiServer.exampleHTTPEndpoint = NewExampleHTTPEndpoint(newApiServer)

	// Initialize the HTTP endpoint server
	// that's hosting the underlying webwire server
	newApiServer.httpSrv = &http.Server{
		Addr:    conf.ServerAddress,
		Handler: newApiServer,
	}

	// Initialize metrics server
	metricsHandler := NewMetricsHandler()
	newApiServer.metricsSrv = &http.Server{
		Addr:    conf.MetricsServerAddress,
		Handler: metricsHandler,
	}

	return newApiServer, nil
}
