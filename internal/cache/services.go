package cache

import (
	"sync"
	"time"

	"github.com/martengine/reception"
)

const cacheTTL = 5 * time.Minute

var c *cache

type cache struct {
	services map[string]reception.Service
	mutex    sync.RWMutex
}

func init() {
	c = &cache{}
	ticker := time.NewTicker(cacheTTL)
	for {
		saveCache(fetchServices())

		<-ticker.C
	}
}

func fetchServices() []reception.Service {
	return []reception.Service{}
}

func saveCache(services []reception.Service) {
	for _, service := range services {
		c.mutex.Lock()
		c.services[service.Name] = service
		c.mutex.Unlock()
	}
}

// Service returns known service instance. If false is returned - service is unknown.
func Service(name string) (reception.Service, bool) {
	c.mutex.RLock()
	service, ok := c.services[name]
	c.mutex.RUnlock()
	return service, ok
}
