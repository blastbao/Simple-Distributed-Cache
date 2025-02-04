package cache

import "sync"

type inMemoryCache struct {
	cacheMap map[string][]byte
	mutex    sync.RWMutex
	Stat
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp, exist := c.cacheMap[k]
	if exist {
		c.del(k, tmp)
	}
	c.cacheMap[k] = v
	c.add(k, v)
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cacheMap[k], nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.cacheMap[k]
	if exist {
		delete(c.cacheMap, k)
		c.del(k, v)
	}
	return nil
}

func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func (c *inMemoryCache) NewScanner() Scanner {

	pairCh := make(chan *pair)

	closeCh := make(chan struct{})

	go func() {
		defer close(pairCh)
		c.mutex.RLock()
		for k, v := range c.cacheMap {
			c.mutex.RUnlock()
			select {
			case <-closeCh:
				return
			case pairCh <- &pair{k, v}:
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}()

	return &inMemoryScanner{pair{}, pairCh, closeCh}
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
}
