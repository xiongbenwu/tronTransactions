package server

type Client interface {
	GetTransactions(address string)
}
