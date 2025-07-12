package metrics

import (
	"sync"
	"time"
)

// Metrics provides basic application metrics
type Metrics struct {
	mu sync.RWMutex

	// HTTP metrics
	requestCount   map[string]int64
	requestLatency map[string][]time.Duration
	errorCount     map[string]int64

	// Database metrics
	dbQueryCount   int64
	dbQueryLatency []time.Duration
	dbErrorCount   int64

	// Application metrics
	activeUsers int64
	uptime      time.Time
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		requestCount:   make(map[string]int64),
		requestLatency: make(map[string][]time.Duration),
		errorCount:     make(map[string]int64),
		uptime:         time.Now(),
	}
}

// RecordRequest records an HTTP request
func (m *Metrics) RecordRequest(method, path string, duration time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := method + " " + path
	m.requestCount[key]++

	if len(m.requestLatency[key]) < 1000 { // Keep last 1000 requests
		m.requestLatency[key] = append(m.requestLatency[key], duration)
	}

	if isError {
		m.errorCount[key]++
	}
}

// RecordDBQuery records a database query
func (m *Metrics) RecordDBQuery(duration time.Duration, isError bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.dbQueryCount++

	if len(m.dbQueryLatency) < 1000 { // Keep last 1000 queries
		m.dbQueryLatency = append(m.dbQueryLatency, duration)
	}

	if isError {
		m.dbErrorCount++
	}
}

// SetActiveUsers sets the number of active users
func (m *Metrics) SetActiveUsers(count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeUsers = count
}

// GetStats returns current metrics statistics
func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := map[string]interface{}{
		"uptime":        time.Since(m.uptime).String(),
		"active_users":  m.activeUsers,
		"db_queries":    m.dbQueryCount,
		"db_errors":     m.dbErrorCount,
		"request_stats": m.requestCount,
		"error_stats":   m.errorCount,
	}

	// Calculate average latencies
	if len(m.dbQueryLatency) > 0 {
		var total time.Duration
		for _, d := range m.dbQueryLatency {
			total += d
		}
		stats["avg_db_latency"] = total / time.Duration(len(m.dbQueryLatency))
	}

	return stats
}

// Reset resets all metrics
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestCount = make(map[string]int64)
	m.requestLatency = make(map[string][]time.Duration)
	m.errorCount = make(map[string]int64)
	m.dbQueryCount = 0
	m.dbQueryLatency = nil
	m.dbErrorCount = 0
	m.activeUsers = 0
	m.uptime = time.Now()
}
