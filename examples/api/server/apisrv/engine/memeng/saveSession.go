package memeng

import wwr "github.com/qbeon/webwire-go"

// SaveSession implements the Engine interface
func (eng *engine) SaveSession(newSession *wwr.Session) error {
	eng.lock.Lock()
	defer eng.lock.Unlock()

	eng.sessions[newSession.Key] = newSession.Clone()
	return nil
}
