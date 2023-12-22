package implement

import (
	"network"
	"sync"
)

type BaseSessionManager[k comparable, v network.ISession] struct {
	sessions map[k]v
	mu       sync.RWMutex
}

func (b *BaseSessionManager[k, v]) AddSession(id k, session v) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if _, exist := b.sessions[id]; exist {
		return false
	}
	b.sessions[id] = session
	return true
}

func (b *BaseSessionManager[k, v]) DelSession(id k) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	delete(b.sessions, id)
}

func (b *BaseSessionManager[k, v]) GeSession(id k) v {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.sessions[id]
}

func (b *BaseSessionManager[k, v]) CountSessions() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.sessions)
}

func (b *BaseSessionManager[k, v]) Clear() {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for id := range b.sessions {
		delete(b.sessions, id)
	}
}

func (b *BaseSessionManager[k, v]) RangeSessions(f func(id k, session v) bool) {
	for id, session := range b.sessions {
		if f(id, session) {
			break
		}
	}
}
