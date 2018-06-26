package client

import (
	"sync/atomic"
	"time"

	webwire "github.com/qbeon/webwire-go"
)

// tryAutoconnect tries to connect to the server.
// If autoconnect is enabled it will spawn a new autoconnector goroutine which
// will periodically poll the server and check whether it's available again.
// If the autoconnector goroutine has already been spawned then it'll
// just await the connection or timeout respectively blocking the calling
// goroutine
func (clt *Client) tryAutoconnect(timeout time.Duration) error {
	autoconn := atomic.LoadInt32(&clt.autoconnect)
	if atomic.LoadInt32(&clt.status) == StatConnected {
		return nil
	} else if autoconn != autoconnectEnabled {
		// Don't try to auto-connect if it's either temporarily deactivated
		// or completely disabled
		return webwire.DisconnectedErr{}
	}

	// Start the reconnector goroutine if not already started.
	// If it's already started then just proceed to wait until either connected or timed out
	clt.backgroundReconnect()

	if timeout > 0 {
		// Await with timeout
		return clt.backReconn.await(timeout)
	}
	// Await indefinitely
	return clt.backReconn.await(0)
}
