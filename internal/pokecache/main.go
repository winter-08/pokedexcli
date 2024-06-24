package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
  cache map[string]cacheEntry
  mu sync.Mutex
  //interval time.Duration
}

type cacheEntry struct {
  createdAt time.Time
  val []byte
}

func (self *Cache) Add(key string, val []byte) {
  self.cache[key] = cacheEntry {
    createdAt: time.Now(),
    val: val,
  }
}

func NewCache(interval time.Duration) (*Cache, error) {
  cache := Cache {
    cache: make(map[string]cacheEntry),
    //interval: interval,
  }

  //defer cache.reapLoop(interval)

  go func ()  {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
      select {
        case <-ticker.C:
          cache.reapLoop(interval)
      }
    }
  }()

  return &cache, nil
}

func (self *Cache) Get(key string) (*[]byte, bool){
  x, found := self.cache[key]
  return &x.val, found 
}

func (self *Cache) reapLoop(interval time.Duration) {
  for key, value := range self.cache {
    if interval >  time.Now().Sub(value.createdAt) {
      delete(self.cache, key)
    }
  }
}
