package goholidays

import (
	"sync"
	"time"
)

// HolidayPool manages reusable Holiday objects to reduce memory allocations
type HolidayPool struct {
	pool sync.Pool
}

// GlobalHolidayPool is a shared pool for Holiday objects
var GlobalHolidayPool = &HolidayPool{
	pool: sync.Pool{
		New: func() interface{} {
			return &Holiday{}
		},
	},
}

// Get retrieves a Holiday object from the pool
func (p *HolidayPool) Get() *Holiday {
	return p.pool.Get().(*Holiday)
}

// Put returns a Holiday object to the pool
func (p *HolidayPool) Put(h *Holiday) {
	// Reset the Holiday object before putting it back
	h.Name = ""
	h.Date = time.Time{}
	h.Category = ""
	h.Observed = nil
	h.Languages = nil
	h.IsObserved = false
	
	p.pool.Put(h)
}

// StringInterner reduces memory usage by interning commonly used strings
type StringInterner struct {
	mu      sync.RWMutex
	strings map[string]string
}

// GlobalStringInterner is a shared string interner
var GlobalStringInterner = &StringInterner{
	strings: make(map[string]string),
}

// Intern returns the canonical version of a string, reducing memory usage
func (si *StringInterner) Intern(s string) string {
	si.mu.RLock()
	if interned, exists := si.strings[s]; exists {
		si.mu.RUnlock()
		return interned
	}
	si.mu.RUnlock()
	
	si.mu.Lock()
	defer si.mu.Unlock()
	
	// Double-check after acquiring write lock
	if interned, exists := si.strings[s]; exists {
		return interned
	}
	
	si.strings[s] = s
	return s
}

// ClearCache clears the string cache (useful for testing or memory management)
func (si *StringInterner) ClearCache() {
	si.mu.Lock()
	defer si.mu.Unlock()
	si.strings = make(map[string]string)
}

// GetCacheSize returns the number of interned strings
func (si *StringInterner) GetCacheSize() int {
	si.mu.RLock()
	defer si.mu.RUnlock()
	return len(si.strings)
}

// OptimizedHoliday creates a memory-optimized Holiday with interned strings
func OptimizedHoliday(name string, date time.Time, category HolidayCategory, languages map[string]string) *Holiday {
	h := GlobalHolidayPool.Get()
	
	h.Name = GlobalStringInterner.Intern(name)
	h.Date = date
	h.Category = category
	h.IsObserved = false
	
	if languages != nil {
		h.Languages = make(map[string]string, len(languages))
		for lang, translation := range languages {
			h.Languages[GlobalStringInterner.Intern(lang)] = GlobalStringInterner.Intern(translation)
		}
	}
	
	return h
}

// BatchDateNormalization optimizes date normalization for multiple dates
func BatchDateNormalization(dates []time.Time) []time.Time {
	normalized := make([]time.Time, len(dates))
	for i, date := range dates {
		normalized[i] = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	}
	return normalized
}

// HolidayCache provides an LRU-style cache for computed holidays
type HolidayCache struct {
	mu       sync.RWMutex
	cache    map[string]map[time.Time]*Holiday
	maxSize  int
	accessed map[string]time.Time
}

// NewHolidayCache creates a new holiday cache with specified max size
func NewHolidayCache(maxSize int) *HolidayCache {
	return &HolidayCache{
		cache:    make(map[string]map[time.Time]*Holiday),
		maxSize:  maxSize,
		accessed: make(map[string]time.Time),
	}
}

// Get retrieves holidays from cache
func (hc *HolidayCache) Get(key string) (map[time.Time]*Holiday, bool) {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	holidays, exists := hc.cache[key]
	if exists {
		hc.accessed[key] = time.Now()
	}
	return holidays, exists
}

// Set stores holidays in cache with LRU eviction
func (hc *HolidayCache) Set(key string, holidays map[time.Time]*Holiday) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	// Evict oldest if at capacity
	if len(hc.cache) >= hc.maxSize {
		oldestKey := ""
		oldestTime := time.Now()
		
		for k, accessTime := range hc.accessed {
			if accessTime.Before(oldestTime) {
				oldestTime = accessTime
				oldestKey = k
			}
		}
		
		if oldestKey != "" {
			delete(hc.cache, oldestKey)
			delete(hc.accessed, oldestKey)
		}
	}
	
	hc.cache[key] = holidays
	hc.accessed[key] = time.Now()
}

// Clear clears the cache
func (hc *HolidayCache) Clear() {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	hc.cache = make(map[string]map[time.Time]*Holiday)
	hc.accessed = make(map[string]time.Time)
}

// Size returns the current cache size
func (hc *HolidayCache) Size() int {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return len(hc.cache)
}
