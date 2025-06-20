package client

import (
	"context"
	"log"
	"math/rand"
	"sync"

	"github.com/adrianpk/rida/internal/cfg"
)

const (
	ottawaLatMin, ottawaLatMax     = 45.40, 45.44
	ottawaLngMin, ottawaLngMax     = -75.72, -75.68
	montrealLatMin, montrealLatMax = 45.49, 45.52
	montrealLngMin, montrealLngMax = -73.59, -73.55
)

type SimManager struct {
	Sims []*Sim
}

func NewClientManager(config *cfg.Config) *SimManager {
	var sims []*Sim

	// Ottawa clients
	for i := 0; i < config.Clients.OttawaQty; i++ {
		lat := ottawaLatMin + rand.Float64()*(ottawaLatMax-ottawaLatMin)
		lng := ottawaLngMin + rand.Float64()*(ottawaLngMax-ottawaLngMin)
		sims = append(sims, NewSim(config.APIKey, lat, lng, "ottawa"))
	}

	// Montreal clients
	for i := 0; i < config.Clients.MontrealQty; i++ {
		lat := montrealLatMin + rand.Float64()*(montrealLatMax-montrealLatMin)
		lng := montrealLngMin + rand.Float64()*(montrealLngMax-montrealLngMin)
		sims = append(sims, NewSim(config.APIKey, lat, lng, "montreal"))
	}

	return &SimManager{Sims: sims}
}

func (m *SimManager) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for _, sim := range m.Sims {
		wg.Add(1)
		go func(s *Sim) {
			defer wg.Done()
			err := s.Run(ctx)
			if err != nil {
				log.Printf("[%s] - Error: %v", s.TagID(), err)
			}
		}(sim)
	}
	wg.Wait()
}
