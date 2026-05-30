package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"predi-backend/internal/ledger"
	"predi-backend/internal/market"
	"predi-backend/internal/streaming"
	"predi-backend/pkg/models"

	_ "github.com/lib/pq"
)

type Server struct {
	db            *sql.DB
	ledgerService *ledger.Service
	marketService *market.Service
	lkService     *streaming.LiveKitService
	users         map[string]*models.User
}

func main() {
	lkAPIKey := os.Getenv("LIVEKIT_API_KEY")
	lkAPISecret := os.Getenv("LIVEKIT_API_SECRET")
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://predi:predi_password@localhost:5432/predi?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := &Server{
		db:            db,
		ledgerService: ledger.NewService(db),
		marketService: market.NewService(db),
		lkService:     streaming.NewLiveKitService(lkAPIKey, lkAPISecret),
		users:         make(map[string]*models.User),
	}

	setupDB(db)

	s.users["user1"] = &models.User{ID: "user1", Username: "alice", Role: "viewer", TrustScore: 85}
	s.users["system"] = &models.User{ID: "system", Username: "platform_seed", Role: "admin"}
	
	s.ledgerService.Deposit("user1", 1000)
	s.ledgerService.Deposit("system", 0)

	s.marketService.CreateMarket(&models.Market{
		ID:                 "00000000-0000-0000-0000-000000000123",
		Title:              "Will 20 or more cars cross the line in the next 60 seconds?",
		Status:             models.StatusOpen,
		EventWindowSeconds: 60,
		ResolutionMethod:   "traffic_count",
		Rules:              `{"targetCount": 20}`,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/ledger/balance", s.handleBalance)
	mux.HandleFunc("/api/ledger/audit", s.handleLedgerAudit)
	mux.HandleFunc("/api/markets", s.handleListMarkets)
	mux.HandleFunc("/api/users/risk", s.handleUserRisk)
	mux.HandleFunc("/api/streaming/token", s.handleGetStreamToken)
	mux.HandleFunc("/api/admin/trigger-resolution", s.handleTriggerResolution)
	mux.HandleFunc("/api/admin/market/health", s.handleMarketHealthUpdate)
	mux.HandleFunc("/api/admin/market/moderate", s.handleMarketModerate)
	mux.HandleFunc("/internal/ai/resolution", s.handleAIResolution)

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	fmt.Println("Predi Go Backend starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler(mux)))
}

func setupDB(db *sql.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS wallets (user_id TEXT PRIMARY KEY, balance DECIMAL(18,2) DEFAULT 0)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS ledger_entries (id SERIAL PRIMARY KEY, user_id TEXT, amount DECIMAL(18,2), type TEXT, reference_id TEXT, created_at TIMESTAMP)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS markets (id TEXT PRIMARY KEY, title TEXT, status TEXT, event_window_seconds INT, resolution_method TEXT, rules TEXT)`)
}

func (s *Server) handleMarketHealthUpdate(w http.ResponseWriter, r *http.Request) {
	marketID := r.URL.Query().Get("market_id")
	health := r.URL.Query().Get("health")
	s.marketService.UpdateHealth(marketID, health, []string{})
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (s *Server) handleMarketModerate(w http.ResponseWriter, r *http.Request) {
	marketID := r.URL.Query().Get("market_id")
	action := r.URL.Query().Get("action")
	status := models.StatusOpen
	if action == "reject" {
		status = models.StatusVoid
	}
	s.marketService.TransitionStatus(marketID, status)
	json.NewEncoder(w).Encode(map[string]string{"status": string(status)})
}

func (s *Server) handleTriggerResolution(w http.ResponseWriter, r *http.Request) {
	marketID := r.URL.Query().Get("market_id")
	market, ok := s.marketService.GetMarket(marketID)
	if !ok {
		http.Error(w, "Market not found", http.StatusNotFound)
		return
	}
	s.marketService.TransitionStatus(marketID, models.StatusResolutionPending)
	s.triggerAIResolution(marketID, "http://mock-stream/hls/123.m3u8", market.EventWindowSeconds)
	json.NewEncoder(w).Encode(map[string]string{"status": "resolution_triggered"})
}

func (s *Server) handleGetStreamToken(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	identity := r.URL.Query().Get("identity")
	token, err := s.lkService.CreateJoinToken(room, identity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *Server) handleListMarkets(w http.ResponseWriter, r *http.Request) {
	markets := s.marketService.ListMarkets()
	json.NewEncoder(w).Encode(markets)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) handleBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	balance := s.ledgerService.GetBalance(userID)
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

func (s *Server) handleLedgerAudit(w http.ResponseWriter, r *http.Request) {
	balance := s.ledgerService.GetBalance("system")
	audit := map[string]interface{}{"total_system_balance": balance}
	json.NewEncoder(w).Encode(audit)
}

func (s *Server) handleUserRisk(w http.ResponseWriter, r *http.Request) {
	riskyUsers := make([]*models.User, 0)
	for _, u := range s.users {
		if u.Suspicious {
			riskyUsers = append(riskyUsers, u)
		}
	}
	json.NewEncoder(w).Encode(riskyUsers)
}

func (s *Server) handleAIResolution(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MarketID   string   `json:"market_id"`
		Confidence float64  `json:"confidence"`
		Flags      []string `json:"flags"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	market, ok := s.marketService.GetMarket(req.MarketID)
	if !ok {
		http.Error(w, "Market not found", http.StatusNotFound)
		return
	}

	// Phase 4 Anti-Manipulation: Check for camera movement
	for _, flag := range req.Flags {
		if flag == "CAMERA_MOVEMENT_DETECTED" {
			s.marketService.TransitionStatus(req.MarketID, models.StatusVoid)
			json.NewEncoder(w).Encode(map[string]string{"status": "voided", "reason": "camera manipulation detected"})
			return
		}
	}

	// Resolution Logic
	s.marketService.TransitionStatus(req.MarketID, models.StatusResolved)
	
	// Deduct Platform Fee (5.0 Play money) and move to system wallet
	s.ledgerService.AddEntry("user1", -5.0, "fee", req.MarketID)
	s.ledgerService.AddEntry("system", 5.0, "platform_fee", req.MarketID)

	json.NewEncoder(w).Encode(map[string]string{"status": "resolved"})
}

func (s *Server) triggerAIResolution(marketID string, streamURL string, window int) {
	payload := map[string]interface{}{"market_id": marketID, "stream_url": streamURL, "window_seconds": window}
	body, _ := json.Marshal(payload)
	go http.Post("http://localhost:8000/tasks", "application/json", bytes.NewBuffer(body))
}
