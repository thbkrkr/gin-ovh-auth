package ovhauth

import "sync"

type CKeyCache struct {
	lock sync.RWMutex
	Map  map[string]string
}

func (this *CKeyCache) get(key string) (v string, ko bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if v, ok := this.Map[key]; ok {
		return v, false
	}
	return "", true
}

func (this *CKeyCache) set(key string, value string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Map[key] = value
}

func (this *CKeyCache) delete(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.Map, key)
}
