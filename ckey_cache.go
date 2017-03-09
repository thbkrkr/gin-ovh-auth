package ovhauth

import "sync"

// cKeyCache stores consumer key during the time of the validation
// with the URL to redirect user in order to log in
type cKeyCache struct {
	lock sync.RWMutex
	Map  map[string]string
}

func (cache *cKeyCache) get(key string) (v string, ko bool) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	if v, ok := cache.Map[key]; ok {
		return v, false
	}
	return "", true
}

func (cache *cKeyCache) set(key string, value string) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	cache.Map[key] = value
}

func (cache *cKeyCache) delete(key string) {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	delete(cache.Map, key)
}
