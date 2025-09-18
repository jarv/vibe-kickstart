package main

import (
	"context"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// ConnectionManager keeps track of active websocket connections
type ConnectionManager struct {
	// Map of name to ConnectionInfo
	connections map[string]ConnectionInfo
	mutex       sync.RWMutex
}

// ConnectionInfo holds information about connections for a specific name
type ConnectionInfo struct {
	Conns []*websocket.Conn
	Name  string
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]ConnectionInfo),
	}
}

// Add adds a connection to the manager
func (cm *ConnectionManager) Add(name string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Get existing info or create new one
	info, exists := cm.connections[name]
	if !exists {
		info = ConnectionInfo{
			Conns: []*websocket.Conn{},
			Name:  name,
		}
	}

	// Append the connection to the list
	info.Conns = append(info.Conns, conn)
	cm.connections[name] = info
}

// Remove removes a connection from the manager
func (cm *ConnectionManager) Remove(name string, conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	info, exists := cm.connections[name]
	if !exists {
		return // Name doesn't exist, nothing to remove
	}

	// Find and remove the connection from the list
	for i, existingConn := range info.Conns {
		if existingConn == conn {
			// Remove the connection at index i
			info.Conns = slices.Delete(info.Conns, i, i+1)
			break
		}
	}

	// If no connections left, remove the name entirely
	if len(info.Conns) == 0 {
		delete(cm.connections, name)
	} else {
		// Update the connection info
		cm.connections[name] = info
	}
}

// BroadcastAll sends a message to all connected clients with timeout protection
func (cm *ConnectionManager) BroadcastAll(ctx context.Context, message []byte) {
	cm.mutex.RLock()

	// Create a copy of connections to release lock quickly
	var allConns []*websocket.Conn
	var allNames []string

	for _, info := range cm.connections {
		for _, conn := range info.Conns {
			allConns = append(allConns, conn)
			allNames = append(allNames, info.Name)
		}
	}
	cm.mutex.RUnlock()

	// Broadcast to all connections concurrently with timeout
	var wg sync.WaitGroup
	for i, conn := range allConns {
		wg.Add(1)
		go func(conn *websocket.Conn, name string) {
			defer wg.Done()

			// Create timeout context for this write (5 second timeout)
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err := conn.Write(writeCtx, websocket.MessageText, message)
			if err != nil {
				log.Printf("Error broadcasting to client %s: %v", name, err)
				// Failed connections will be cleaned up by the handler goroutine
			}
		}(conn, allNames[i])
	}

	// Wait for all writes to complete or timeout
	wg.Wait()
}

// BroadcastToName sends a message to all connections with a specific name
func (cm *ConnectionManager) BroadcastToName(ctx context.Context, name string, message []byte) {
	cm.mutex.RLock()

	info, exists := cm.connections[name]
	if !exists {
		cm.mutex.RUnlock()
		return // Name doesn't exist
	}

	// Create a copy of connections to release lock quickly
	connsCopy := make([]*websocket.Conn, len(info.Conns))
	copy(connsCopy, info.Conns)
	cm.mutex.RUnlock()

	// Broadcast to connections concurrently with timeout
	var wg sync.WaitGroup
	for _, conn := range connsCopy {
		wg.Add(1)
		go func(conn *websocket.Conn) {
			defer wg.Done()

			// Create timeout context for this write (5 second timeout)
			writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err := conn.Write(writeCtx, websocket.MessageText, message)
			if err != nil {
				log.Printf("Error broadcasting to client %s: %v", name, err)
				// Failed connections will be cleaned up separately
			}
		}(conn)
	}

	// Wait for all writes to complete or timeout
	wg.Wait()
}

// GetConnectionCount returns the total number of active connections
func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	count := 0
	for _, info := range cm.connections {
		count += len(info.Conns)
	}
	return count
}

// GetNameCount returns the number of unique names with connections
func (cm *ConnectionManager) GetNameCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.connections)
}

// GetConnectionsByName returns the number of connections for a specific name
func (cm *ConnectionManager) GetConnectionsByName(name string) int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	info, exists := cm.connections[name]
	if !exists {
		return 0
	}
	return len(info.Conns)
}

// HasName checks to see if a player's name is present in the connection list
func (cm *ConnectionManager) HasName(name string) bool {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	_, exists := cm.connections[name]

	return exists
}
