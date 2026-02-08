package session

import (
	"sync"
	"time"
)

// Session tracks movies that have been used in a game session
type Session struct {
	ID           string
	UsedMovieIDs map[int64]bool
	CreatedAt    time.Time
	LastAccessed time.Time
}

// Manager handles session storage and operations
type Manager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	timeout  time.Duration
}

// NewManager creates a new session manager with specified timeout
func NewManager(timeout time.Duration) *Manager {
	m := &Manager{
		sessions: make(map[string]*Session),
		timeout:  timeout,
	}
	go m.cleanupLoop()
	return m
}

// CreateSession creates a new session with the given ID
func (m *Manager) CreateSession(sessionID string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.sessions[sessionID]; exists {
		return m.sessions[sessionID]
	}

	m.sessions[sessionID] = &Session{
		ID:           sessionID,
		UsedMovieIDs: make(map[int64]bool),
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}
	return m.sessions[sessionID]
}

// GetSession retrieves a session by ID
func (m *Manager) GetSession(sessionID string) *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, exists := m.sessions[sessionID]
	if exists {
		session.LastAccessed = time.Now()
	}
	return session
}

// MarkMovieAsUsed marks a movie ID as used for the session
func (m *Manager) MarkMovieAsUsed(sessionID string, movieID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, exists := m.sessions[sessionID]; exists {
		session.UsedMovieIDs[movieID] = true
		session.LastAccessed = time.Now()
	}
}

// IsMovieUsed checks if a movie has been used in the session
func (m *Manager) IsMovieUsed(sessionID string, movieID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if session, exists := m.sessions[sessionID]; exists {
		return session.UsedMovieIDs[movieID]
	}
	return false
}

// DeleteSession removes a session
func (m *Manager) DeleteSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, sessionID)
}

// cleanupLoop periodically removes expired sessions
func (m *Manager) cleanupLoop() {
	for {
		time.Sleep(m.timeout / 2)
		m.mu.Lock()
		now := time.Now()
		for id, session := range m.sessions {
			if now.Sub(session.LastAccessed) > m.timeout {
				delete(m.sessions, id)
			}
		}
		m.mu.Unlock()
	}
}
