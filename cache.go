package reception

import (
	"sync"
	"time"

	"github.com/martengine/reception/service"
)

const cacheTTL = 5 * time.Minute

var c *cache

type cache struct {
	services map[string]service.Service
	mutex    sync.RWMutex
}

func init() {
	c = &cache{}

	go func(ticker *time.Ticker) {
		// refresh cache from time to time.
		for {
			saveCache(fetchServices())

			<-ticker.C
		}
	}(time.NewTicker(cacheTTL))
}

func fetchServices() []service.Service {
	return []service.Service{}
}

func saveCache(services []service.Service) {
	for _, service := range services {
		c.mutex.Lock()
		c.services[service.Name] = service
		c.mutex.Unlock()
	}
}

// ServiceByName returns known service instance. If false is returned - service is unknown.
func ServiceByName(name string) (service.Service, bool) {
	c.mutex.RLock()
	service, ok := c.services[name]
	c.mutex.RUnlock()
	return service, ok
}

// PublicServices returns a list of known public services.
func PublicServices() []service.Service {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var services []service.Service
	for _, service := range c.services {
		if !service.Public {
			continue
		}

		services = append(services, service)
	}

	return services
}
