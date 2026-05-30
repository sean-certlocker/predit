package ledger

import (
	"database/sql"
	"errors"
	"predi-backend/pkg/models"
	"time"

	_ "github.com/lib/pq"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetBalance(userID string) float64 {
	var balance float64
	err := s.db.QueryRow("SELECT balance FROM wallets WHERE user_id = $1", userID).Scan(&balance)
	if err != nil {
		return 0
	}
	return balance
}

func (s *Service) AddEntry(userID string, amount float64, entryType string, refID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update or create wallet
	var currentBalance float64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE", userID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create wallet if it doesn't exist (for MVP convenience)
			if amount < 0 {
				return errors.New("insufficient balance")
			}
			_, err = tx.Exec("INSERT INTO wallets (user_id, balance) VALUES ($1, $2)", userID, amount)
		} else {
			return err
		}
	} else {
		if currentBalance+amount < 0 {
			return errors.New("insufficient balance")
		}
		_, err = tx.Exec("UPDATE wallets SET balance = balance + $1 WHERE user_id = $2", amount, userID)
	}

	if err != nil {
		return err
	}

	// Add ledger entry
	_, err = tx.Exec("INSERT INTO ledger_entries (user_id, amount, type, reference_id, created_at) VALUES ($1, $2, $3, $4, $5)",
		userID, amount, entryType, nil, time.Now()) // refID skipped for simple MVP UUID compatibility

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) Deposit(userID string, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid deposit amount")
	}
	return s.AddEntry(userID, amount, "deposit", "")
}
