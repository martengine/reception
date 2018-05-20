package service

import (
	"sync"
	"time"
)

const cacheTTL = 5 * time.Minute

var c *cache

type cache struct {
	services map[string]Service
	mutex    sync.RWMutex
}

func init() {
	c = &cache{}

	go func(ticker *time.Ticker) {
		// refresh cache from time to time.
		for {
			saveCache(fetch())

			<-ticker.C
		}
	}(time.NewTicker(cacheTTL))
}

func fetch() []Service {
	return []Service{}
}

func saveCache(services []Service) {
	for _, service := range services {
		c.mutex.Lock()
		c.services[service.Name] = service
		c.mutex.Unlock()
	}
}

// ByName returns known service instance. If false is returned - service is unknown.
func ByName(name string) (Service, bool) {
	c.mutex.RLock()
	service, ok := c.services[name]
	c.mutex.RUnlock()
	return service, ok
}

// Public returns a list of known public services.
func Public() []Service {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var services []Service
	for _, service := range c.services {
		if !service.Public {
			continue
		}

		services = append(services, service)
	}

	return services
}
