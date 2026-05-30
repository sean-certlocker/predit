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
	"sync"
)

type Server struct {
	ledgerService *ledger.Service
	marketService *market.Service
	lkService     *streaming.LiveKitService
}

func main() {
	lkAPIKey := os.Getenv("LIVEKIT_API_KEY")
	lkAPISecret := os.Getenv("LIVEKIT_API_SECRET")

	s := &Server{
		ledgerService: ledger.NewService(),
		marketService: market.NewService(),
		lkService:     streaming.NewLiveKitService(lkAPIKey, lkAPISecret),
	}

	// Seed a test market
	s.marketService.CreateMarket(&models.Market{
		ID:                 "123",
		Title:              "Will 20 or more cars cross the line in the next 60 seconds?",
		Status:             models.StatusOpen,
		EventWindowSeconds: 60,
		ResolutionMethod:   "traffic_count",
		Rules:              `{"targetCount": 20}`,
	})

	s.ledgerService.Deposit("user1", 1000)

	// Routes
	http.HandleFunc("/api/health", s.handleHealth)
	http.HandleFunc("/api/ledger/balance", s.handleBalance)
	http.HandleFunc("/api/markets", s.handleListMarkets)
	http.HandleFunc("/api/streaming/token", s.handleGetStreamToken)
	http.HandleFunc("/api/admin/trigger-resolution", s.handleTriggerResolution)
	http.HandleFunc("/internal/ai/resolution", s.handleAIResolution)

	fmt.Println("Predi Go Backend starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
