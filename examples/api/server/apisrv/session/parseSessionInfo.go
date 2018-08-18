package session

import (
	"github.com/pkg/errors"
	wwr "github.com/qbeon/webwire-go"
)

// ParseSessionInfo parses the session info from a variant map
func ParseSessionInfo(data map[string]interface{}) wwr.SessionInfo {
	// Parse identifier
	var userIdentifier string

	switch val := data["id"].(type) {
	case string:
		userIdentifier = val
	}

	// Parse client type
	var userType UserType

	switch val := data["type"].(type) {
	case UserType:
		userType = val
	case string:
		if err := userType.FromString(val); err != nil {
			panic(errors.Wrap(
				err,
				"couldn't parse UserType from session info",
			))
		}
	}

	return &SessionInfo{
		UserIdentifier: userIdentifier,
		UserType:       userType,
	}
}
