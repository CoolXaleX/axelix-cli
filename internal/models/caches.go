package models

// CachesFeed is the response of the caches actuator endpoint.
type CachesFeed struct {
	CacheManagers []CacheManager `json:"cacheManagers"`
}

// CacheManager holds a named set of caches.
type CacheManager struct {
	Name   string  `json:"name"`
	Caches []Cache `json:"caches"`
}

// Cache represents a single cache inside a manager.
type Cache struct {
	Name         string `json:"name"`
	Target       string `json:"target"`
	Enabled      bool   `json:"enabled"`
	ContainsStats bool  `json:"containsStats"`
}
