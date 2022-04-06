package server

import (
	"encoding/json"
	"net/http"
	"transactions/pkg/tron"

	"github.com/julienschmidt/httprouter"
)

type Client interface {
	GetTransactions(address string) ([]tron.Transaction, error)
}

type Server struct {
	Client Client
}

func (s *Server) GetTransactions(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	address := ps.ByName("tx_id")
	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	transactions, err := s.Client.GetTransactions(address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	bytes, err := json.Marshal(transactions)
	if err != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
}

func (s *Server) ListenAndServe() error {
	router := httprouter.New()
	router.GET("/transactions/:tx_id", s.GetTransactions)

	srv := http.Server{Addr: ":8080", Handler: router}

	return srv.ListenAndServe()
}
