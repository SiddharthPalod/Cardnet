package registry

import "sync"

type Registry struct {
	mu     sync.RWMutex
	models map[string]interface{}
}

func NewRegistry() *Registry {
	return &Registry{
		models: make(map[string]interface{}),
	}
}

func (r *Registry) Register(name string, model interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.models[name] = model
}

func (r *Registry) Get(name string) interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.models[name]
}
