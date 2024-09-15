package cache

import (
	"errors"
	"sync"
	"time"
)

// defaultExpiration - tiempo por defecto de expiración de los datos en el caché. (una semana)
const defaultExpiration = 7 * 24 * time.Hour

type Cache struct {
	data       map[string]cacheItem
	mutex      sync.RWMutex
	expiration time.Duration
}

type cacheItem struct {
	value     interface{}
	timestamp time.Time
}

// New crea un nuevo caché con el tiempo de expiración especificado.
func New(expiration time.Duration) *Cache {
	if expiration == 0 {
		expiration = defaultExpiration
	}
	return &Cache{
		data:       make(map[string]cacheItem),
		expiration: expiration,
	}
}

// Get devuelve el valor por clave y un indicador de éxito.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, exists := c.data[key]
	if !exists || time.Since(item.timestamp) > c.expiration {
		return nil, false
	}
	return item.value, true
}

// Add agrega un nuevo valor por clave. Devuelve un error si la clave ya existe.
func (c *Cache) Add(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.data[key]; exists {
		return errors.New("key already exists")
	}
	c.data[key] = cacheItem{
		value:     value,
		timestamp: time.Now(),
	}
	return nil
}

// Update modifica el valor por clave. Devuelve un error si la clave no existe.
func (c *Cache) Update(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, exists := c.data[key]; !exists {
		return errors.New("key does not exist")
	}
	c.data[key] = cacheItem{
		value:     value,
		timestamp: time.Now(),
	}
	return nil
}

// Check devuelve un indicador de existencia de un valor por clave.
func (c *Cache) Check(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, exists := c.data[key]
	return exists
}

// Delete elimina un valor por clave.
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

// cleanup elimina los datos caducados del caché.
func (c *Cache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, item := range c.data {
		if time.Since(item.timestamp) > c.expiration {
			delete(c.data, key)
		}
	}
}
