package ledger

import (
	"errors"
	"predi-backend/pkg/models"
	"sync"
	"time"
)

type Service struct {
	mu      sync.RWMutex
	entries []models.LedgerEntry
	wallets map[string]float64
}

func NewService() *Service {
	return &Service{
		entries: make([]models.LedgerEntry, 0),
		wallets: make(map[string]float64),
	}
}

func (s *Service) GetBalance(userID string) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.wallets[userID]
}

func (s *Service) AddEntry(userID string, amount float64, entryType string, refID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for sufficient balance if it's a stake
	if amount < 0 && s.wallets[userID]+amount < 0 {
		return errors.New("insufficient balance")
	}

	entry := models.LedgerEntry{
		ID:          "L-" + time.Now().Format("20060102150405"),
		UserID:      userID,
		Amount:      amount,
		Type:        entryType,
		ReferenceID: refID,
		CreatedAt:   time.Now(),
	}

	s.entries = append(s.entries, entry)
	s.wallets[userID] += amount

	return nil
}

func (s *Service) Deposit(userID string, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid deposit amount")
	}
	return s.AddEntry(userID, amount, "deposit", "")
}
