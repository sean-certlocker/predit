package market

import (
	"errors"
	"predi-backend/pkg/models"
	"sync"
	"time"
)

type Service struct {
	mu      sync.RWMutex
	markets map[string]*models.Market
}

func NewService() *Service {
	return &Service{
		markets: make(map[string]*models.Market),
	}
}

func (s *Service) CreateMarket(market *models.Market) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.markets[market.ID] = market
}

func (s *Service) GetMarket(id string) (*models.Market, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.markets[id]
	return m, ok
}

func (s *Service) ListMarkets() []*models.Market {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*models.Market, 0, len(s.markets))
	for _, m := range s.markets {
		list = append(list, m)
	}
	return list
}

func (s *Service) TransitionStatus(id string, newStatus models.MarketStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.markets[id]
	if !ok {
		return errors.New("market not found")
	}

	if !isValidTransition(m.Status, newStatus) {
		return errors.New("invalid status transition from " + string(m.Status) + " to " + string(newStatus))
	}

	m.Status = newStatus
	return nil
}

func isValidTransition(current, next models.MarketStatus) bool {
	switch current {
	case models.StatusDraft:
		return next == models.StatusOpen
	case models.StatusOpen:
		return next == models.StatusLocked || next == models.StatusVoid
	case models.StatusLocked:
		return next == models.StatusEventActive || next == models.StatusVoid
	case models.StatusEventActive:
		return next == models.StatusResolutionPending || next == models.StatusVoid
	case models.StatusResolutionPending:
		return next == models.StatusResolved || next == models.StatusVoid
	case models.StatusResolved:
		return next == models.StatusSettled
	case models.StatusVoid:
		return next == models.StatusRefunded
	default:
		return false
	}
}
