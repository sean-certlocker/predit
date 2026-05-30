package models

import (
	"time"
)

type MarketStatus string

const (
	StatusDraft             MarketStatus = "DRAFT"
	StatusOpen              MarketStatus = "OPEN"
	StatusLocked            MarketStatus = "LOCKED"
	StatusEventActive       MarketStatus = "EVENT_ACTIVE"
	StatusResolutionPending MarketStatus = "RESOLUTION_PENDING"
	StatusResolved          MarketStatus = "RESOLVED"
	StatusVoid              MarketStatus = "VOID"
	StatusSettled           MarketStatus = "SETTLED"
	StatusRefunded          MarketStatus = "REFUNDED"
)

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	TrustScore int       `json:"trust_score"`
	CreatedAt  time.Time `json:"created_at"`
}

type Market struct {
	ID                           string       `json:"id"`
	Title                        string       `json:"title"`
	Status                       MarketStatus `json:"status"`
	EventWindowSeconds           int          `json:"event_window_seconds"`
	BettingClosesBeforeStartSecs int          `json:"betting_closes_before_start_seconds"`
	ResolutionMethod             string       `json:"resolution_method"`
	MaxStake                     float64      `json:"max_stake"`
	MaxPool                      float64      `json:"max_pool"`
	Rules                        string       `json:"rules"` // JSON string for simplicity in MVP
}

type LedgerEntry struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // stake, win, fee, refund
	ReferenceID string    `json:"reference_id"`
	CreatedAt   time.Time `json:"created_at"`
}
