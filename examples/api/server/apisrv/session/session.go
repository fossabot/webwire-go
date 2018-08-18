package session

import (
	"time"

	wwr "github.com/qbeon/webwire-go"
)

// Session represents an API server session
type Session struct {
	Key        string
	Creation   time.Time
	LastLookup time.Time
	Info       SessionInfo
}

// SessionInfo represents an API server specific session info object
type SessionInfo struct {
	UserIdentifier string
	UserType       UserType
}

// Copy implements the wwr.SessionInfo interface
func (sinf *SessionInfo) Copy() wwr.SessionInfo {
	return &SessionInfo{
		UserIdentifier: sinf.UserIdentifier,
		UserType:       sinf.UserType,
	}
}

// Fields implements the wwr.SessionInfo interface
func (sinf *SessionInfo) Fields() []string {
	return []string{"id", "type"}
}

// Value implements the wwr.SessionInfo interface
func (sinf *SessionInfo) Value(fieldName string) interface{} {
	switch fieldName {
	case "id":
		return sinf.UserIdentifier
	case "type":
		return sinf.UserType.String()
	}
	return nil
}
