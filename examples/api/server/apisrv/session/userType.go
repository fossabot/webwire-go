package session

import (
	"strings"

	"github.com/pkg/errors"
)

// UserType represents the role of an API client
type UserType int

const (
	// Guest represents an unauthenticated guest client
	Guest UserType = iota

	// User represents a regular user client
	User

	// Admin represents an administrator client
	Admin
)

// String stringifies the value
func (tp UserType) String() string {
	switch tp {
	case Guest:
		return "guest"
	case User:
		return "user"
	case Admin:
		return "admin"
	}
	return ""
}

// FromString parses the client type from string
func (tp *UserType) FromString(str string) error {
	switch strings.ToLower(str) {
	case "guest":
		*tp = Guest
	case "user":
		*tp = User
	case "admin":
		*tp = Admin
	default:
		return errors.Errorf(
			"invalid string representation of UserType type: %s",
			str,
		)
	}
	return nil
}
