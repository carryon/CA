package server

import (
	"runtime"
	"sync"
	"time"

	"github.com/manucorporat/stats"
)

type Stats struct {
	sync.RWMutex
	Saved    map[string]uint64
	Messages *stats.StatsCollector
	Users    *stats.StatsCollector
	Ips      *stats.StatsCollector
}

func NewStats() *Stats {
	return &Stats{
		Saved:    make(map[string]uint64),
		Messages: stats.New(),
		Users:    stats.New(),
		Ips:      stats.New(),
	}
}

func (s *Stats) statsWorker() {
	c := time.Tick(1 * time.Second)
	var lastMallocs uint64
	var lastFrees uint64
	for range c {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		s.Lock()
		s.Saved = map[string]uint64{
			"timestamp":  uint64(time.Now().Unix()),
			"HeapInuse":  stats.HeapInuse,
			"StackInuse": stats.StackInuse,
			"Mallocs":    (stats.Mallocs - lastMallocs),
			"Frees":      (stats.Frees - lastFrees),
			"Inbound":    uint64(s.Messages.Get("inbound")),
			"Outbound":   uint64(s.Messages.Get("outbound")),
			"Connected":  s.connectedUsers(),
		}
		lastMallocs = stats.Mallocs
		lastFrees = stats.Frees
		s.Messages.Reset()
		s.Unlock()
	}
}

func (s *Stats) connectedUsers() uint64 {
	connected := s.Users.Get("connected") - s.Users.Get("disconnected")
	if connected < 0 {
		return 0
	}
	return uint64(connected)
}
