package test

import (
	"context"
	"testing"
	"time"

	webwire "github.com/qbeon/webwire-go"
	webwireClient "github.com/qbeon/webwire-go/client"
)

// TestServerInitiatedSessionDestruction verifies
// server-initiated session destruction
func TestServerInitiatedSessionDestruction(t *testing.T) {
	sessionCreationCallbackCalled := newPending(1, 1*time.Second, true)
	sessionDestructionCallbackCalled := newPending(1, 1*time.Second, true)
	var createdSession *webwire.Session
	expectedCredentials := webwire.Payload{
		Encoding: webwire.EncodingUtf8,
		Data:     []byte("secret_credentials"),
	}
	placeholderMessage := webwire.Payload{
		Data: []byte("nothinginteresting"),
	}
	currentStep := 1

	// Initialize webwire server
	server := setupServer(
		t,
		&serverImpl{
			onRequest: func(
				_ context.Context,
				clt *webwire.Client,
				msg *webwire.Message,
			) (webwire.Payload, error) {
				// On step 2 - verify session creation and correctness
				if currentStep == 2 {
					sess := clt.Session()
					compareSessions(t, createdSession, sess)
					if string(msg.Payload.Data) != sess.Key {
						t.Errorf(
							"Clients session key doesn't match: "+
								"client: '%s' | server: '%s'",
							string(msg.Payload.Data),
							sess.Key,
						)
					}
					return webwire.Payload{}, nil
				}

				// on step 3 - close session and verify its destruction
				if currentStep == 3 {
					/******************************************************\
						Server-side session destruction initiation
					\******************************************************/
					// Attempt to destroy this clients session
					// on the end of the first step
					err := clt.CloseSession()
					if err != nil {
						t.Errorf(
							"Couldn't close the active session "+
								"on the server: %s",
							err,
						)
					}

					// Verify destruction
					sess := clt.Session()
					if sess != nil {
						t.Errorf(
							"Expected the session to be destroyed, got: %v",
							sess,
						)
					}

					return webwire.Payload{}, nil
				}

				// On step 4 - verify session destruction
				if currentStep == 4 {
					sess := clt.Session()
					if sess != nil {
						t.Errorf(
							"Expected the session to be destroyed, got: %v",
							sess,
						)
					}
					return webwire.Payload{}, nil
				}

				// On step 1 - authenticate and create a new session
				if err := clt.CreateSession(nil); err != nil {
					return webwire.Payload{}, err
				}

				// Return the key of the newly created session
				return webwire.Payload{
					Data: []byte(clt.SessionKey()),
				}, nil
			},
		},
		webwire.ServerOptions{},
	)

	// Initialize client
	client := newCallbackPoweredClient(
		server.Addr().String(),
		webwireClient.Options{
			DefaultRequestTimeout: 2 * time.Second,
		},
		callbackPoweredClientHooks{
			OnSessionCreated: func(_ *webwire.Session) {
				// Mark the client-side session creation callback executed
				sessionCreationCallbackCalled.Done()
			},
			OnSessionClosed: func() {
				// Ensure this callback is called during the
				if currentStep != 3 {
					t.Errorf(
						"Client-side session destruction callback "+
							"called at wrong step (%d)",
						currentStep,
					)
				}
				sessionDestructionCallbackCalled.Done()
			},
		},
	)

	/*****************************************************************\
		Step 1 - Session Creation
	\*****************************************************************/
	if err := client.connection.Connect(); err != nil {
		t.Fatalf("Couldn't connect: %s", err)
	}

	// Send authentication request
	authReqReply, err := client.connection.Request("login", expectedCredentials)
	if err != nil {
		t.Fatalf("Authentication request failed: %s", err)
	}

	createdSession = client.connection.Session()

	// Verify reply
	comparePayload(
		t,
		"authentication reply",
		webwire.Payload{
			Data: []byte(createdSession.Key),
		},
		authReqReply,
	)

	// Wait for the client-side session creation callback to be executed
	if err := sessionCreationCallbackCalled.Wait(); err != nil {
		t.Fatal("Session creation callback not called")
	}

	// Ensure the session was locally created
	currentSessionAfterCreation := client.connection.Session()
	if currentSessionAfterCreation.Key == "" {
		t.Fatalf(
			"Expected session on client-side, got none: %v",
			currentSessionAfterCreation,
		)
	}

	/*****************************************************************\
		Step 2 - Session Creation Verification
	\*****************************************************************/
	currentStep = 2

	// Send a test-request to verify the session creation on the server
	if _, err := client.connection.Request(
		"",
		webwire.Payload{Data: []byte(client.connection.Session().Key)},
	); err != nil {
		t.Fatalf("Session creation verification request failed: %s", err)
	}

	/*****************************************************************\
		Step 3 - Server-Side Session Destruction
	\*****************************************************************/
	currentStep = 3

	// Request session destruction
	if _, err := client.connection.Request("", placeholderMessage); err != nil {
		t.Fatalf("Session destruction request failed: %s", err)
	}

	// Wait for the client-side session destruction callback to be called
	if err := sessionDestructionCallbackCalled.Wait(); err != nil {
		t.Fatal("Session destruction callback not called")
	}

	/*****************************************************************\
		Step 4 - Destruction Verification
	\*****************************************************************/
	currentStep = 4

	// Ensure the session is destroyed locally as well
	currentSessionAfterDestruction := client.connection.Session()
	if currentSessionAfterDestruction != nil {
		t.Fatalf(
			"Expected session to be destroyed on the client as well, got: %v",
			currentSessionAfterDestruction,
		)
	}

	// Send a test-request to verify the session was destroyed on the server
	if _, err := client.connection.Request("", placeholderMessage); err != nil {
		t.Fatalf("Session destruction verification request failed: %s", err)
	}
}
