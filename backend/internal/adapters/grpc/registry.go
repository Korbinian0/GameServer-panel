package grpcclient

import (
	"sync"
)

type BackendRegistry struct {
	certPath string
	mu       sync.RWMutex
	clients  map[string]*BackendClient
}

func NewRegistry(certPath string) *BackendRegistry {
	return &BackendRegistry{
		certPath: certPath,
		clients:  make(map[string]*BackendClient),
	}
}

func (r *BackendRegistry) AddBackend(nodeID, target string) error {
	client, err := NewBackendClient(target, r.certPath)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if old, ok := r.clients[nodeID]; ok {
		old.Close()
	}
	r.clients[nodeID] = client
	return nil
}

func (r *BackendRegistry) GetBackend(nodeID string) (*BackendClient, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	client, ok := r.clients[nodeID]
	return client, ok
}

func (r *BackendRegistry) RemoveBackend(nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if client, ok := r.clients[nodeID]; ok {
		client.Close()
		delete(r.clients, nodeID)
	}
}

func (r *BackendRegistry) ListBackendIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.clients))
	for id := range r.clients {
		ids = append(ids, id)
	}
	return ids
}

func (r *BackendRegistry) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, client := range r.clients {
		client.Close()
	}
	r.clients = make(map[string]*BackendClient)
}
