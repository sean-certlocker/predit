package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"predi-backend/internal/ledger"
	"predi-backend/internal/market"
	"predi-backend/internal/streaming"
	"predi-backend/pkg/models"
)

type Server struct {
	ledgerService *ledger.Service
	marketService *market.Service
	lkService     *streaming.LiveKitService
	users         map[string]*models.User
}

func main() {
	lkAPIKey := os.Getenv("LIVEKIT_API_KEY")
	lkAPISecret := os.Getenv("LIVEKIT_API_SECRET")

	s := &Server{
		ledgerService: ledger.NewService(),
		marketService: market.NewService(),
		lkService:     streaming.NewLiveKitService(lkAPIKey, lkAPISecret),
		users:         make(map[string]*models.User),
	}

	// Seed data
	s.users["user1"] = &models.User{ID: "user1", Username: "alice", Role: "viewer", TrustScore: 85}
	s.users["user2"] = &models.User{ID: "user2", Username: "bob", Role: "creator", TrustScore: 92}
	s.users["user3"] = &models.User{ID: "user3", Username: "mallory", Role: "viewer", TrustScore: 10, Suspicious: true}

	s.marketService.CreateMarket(&models.Market{
		ID:                 "123",
		Title:              "Will 20 or more cars cross the line in the next 60 seconds?",
		Status:             models.StatusOpen,
		EventWindowSeconds: 60,
		ResolutionMethod:   "traffic_count",
		Rules:              `{"targetCount": 20}`,
		Health:             "Green",
	})

	s.marketService.CreateMarket(&models.Market{
		ID:                 "124",
		Title:              "Will a red truck pass by in the next 30 seconds?",
		Status:             models.StatusDraft,
		EventWindowSeconds: 30,
		ResolutionMethod:   "object_detection",
		Rules:              `{"target": "red_truck"}`,
	})

	s.ledgerService.Deposit("user1", 1000)

	// Routes
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

func (s *Server) handleMarketHealthUpdate(w http.ResponseWriter, r *http.Request) {
	marketID := r.URL.Query().Get("market_id")
	health := r.URL.Query().Get("health") // Green, Yellow, Red
	
	err := s.marketService.UpdateHealth(marketID, health, []string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{"status": "updated", "health": health})
}

func (s *Server) handleMarketModerate(w http.ResponseWriter, r *http.Request) {
	marketID := r.URL.Query().Get("market_id")
	action := r.URL.Query().Get("action") // approve, reject
	
	status := models.StatusOpen
	if action == "reject" {
		status = models.StatusVoid
	}
	
	err := s.marketService.TransitionStatus(marketID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
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
	
	json.NewEncoder(w).Encode(map[string]string{"status": "resolution_triggered", "market_id": marketID})
}

func (s *Server) handleGetStreamToken(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	identity := r.URL.Query().Get("identity")

	if room == "" || identity == "" {
		http.Error(w, "room and identity required", http.StatusBadRequest)
		return
	}

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
	audit := map[string]interface{}{
		"total_system_balance": 1000.0,
		"platform_fees":        25.50,
		"pending_stakes":       450.0,
	}
	json.NewEncoder(w).Encode(audit)
}

func (s *Server) handleUserRisk(w http.ResponseWriter, r *http.Request) {
	riskyUsers := make([]*models.User, 0)
	for _, u := range s.users {
		if u.Suspicious || u.TrustScore < 30 {
			riskyUsers = append(riskyUsers, u)
		}
	}
	json.NewEncoder(w).Encode(riskyUsers)
}

type AIResolutionRequest struct {
	MarketID   string  `json:"market_id"`
	Count      int     `json:"count"`
	Confidence float64 `json:"confidence"`
	Flags      []string `json:"flags"`
}

func (s *Server) handleAIResolution(w http.ResponseWriter, r *http.Request) {
	var req AIResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	market, ok := s.marketService.GetMarket(req.MarketID)
	if !ok {
		http.Error(w, "Market not found", http.StatusNotFound)
		return
	}

	if market.Status != models.StatusResolutionPending {
		http.Error(w, "Market not awaiting resolution", http.StatusBadRequest)
		return
	}

	// Resolution Logic
	if req.Confidence < 0.85 {
		s.marketService.TransitionStatus(req.MarketID, models.StatusVoid)
		// TODO: Call ledgerService to refund
		json.NewEncoder(w).Encode(map[string]string{"status": "voided", "reason": "low confidence"})
		return
	}

	for _, flag := range req.Flags {
		if flag == "NSFW" || flag == "VIOLENCE" {
			s.marketService.TransitionStatus(req.MarketID, models.StatusVoid)
			json.NewEncoder(w).Encode(map[string]string{"status": "voided", "reason": "safety violation: " + flag})
			return
		}
	}

	// Simplified result logic
	s.marketService.TransitionStatus(req.MarketID, models.StatusResolved)
	json.NewEncoder(w).Encode(map[string]string{"status": "resolved", "market_id": req.MarketID})
}

func (s *Server) triggerAIResolution(marketID string, streamURL string, window int) {
	payload := map[string]interface{}{
		"market_id":      marketID,
		"stream_url":     streamURL,
		"window_seconds": window,
	}
	body, _ := json.Marshal(payload)
	
	// Fire and forget or handle error
	go func() {
		http.Post("http://localhost:8000/tasks", "application/json", bytes.NewBuffer(body))
	}()
}
