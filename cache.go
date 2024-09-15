package cache

import (
	"sync"
	"time"
)

// defaultExpiration - tiempo por defecto de expiración de los datos en el caché. (una semana)
const defaultExpiration = 168 * time.Hour

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

// Update Modifica el valor por clave.
func (c *Cache) Update(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheItem{
		value:     value,
		timestamp: time.Now(),
	}
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

// Сохранение refresh токена: Когда пользователь логинится, вы сохраняете refresh токен в кеш с помощью метода Update.

// Проверка и обновление JWT токена: Когда JWT токен истекает, ваше приложение проверяет наличие refresh
// токена в кеше с помощью метода Check и получает его с помощью метода Get. Если refresh токен все еще активен,
// вы используете его для получения нового JWT токена и обновляете кеш.

// Удаление устаревших токенов: Если refresh токен истек (например, прошло 7 дней),
// он автоматически удаляется из кеша, и пользователю придется снова логиниться.
