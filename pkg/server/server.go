package server

import (
	"encoding/json"
	"net/http"
	"transactions/pkg/tron"
)

type Client interface {
	GetTransactions() ([]tron.Transaction, error)
}

type Server struct {
	Client Client
}

func (s *Server) Transactions(w http.ResponseWriter, req *http.Request) {
	transactions, err := s.Client.GetTransactions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(transactions)
	if err != nil {
		return
	}
	w.Write(bytes)
}

func (s *Server) ListenAndServe() {
	http.HandleFunc("/transactions", s.Transactions)
	http.ListenAndServe(":8080", nil)
}
