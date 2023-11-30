package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/thrillee/price-fetcher/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        PriceFetcher
}

func NewJsonAPIServer(listenAddr string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/", makeHTTPHandlerFunc(s.handleFetchPrice))
	log.Println("Server starting...")
	http.ListenAndServe(s.listenAddr, nil)
}

func makeHTTPHandlerFunc(apiFunc APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestID", rand.Intn(100000000))

	return func(w http.ResponseWriter, r *http.Request) {
		if err := apiFunc(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}

func (s *JSONAPIServer) handleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")

	price, err := s.svc.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}

	priceResp := types.PriceResponse{
		Ticker: ticker,
		Price:  price,
	}

	return writeJSON(w, http.StatusOK, *&priceResp)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
