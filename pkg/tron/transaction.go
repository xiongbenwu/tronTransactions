package tron

import "encoding/json"

type Transaction struct {
	ID        string  `json:"txID"`
	Timestamp int64   `json:"block_timestamp"`
	RawData   RawData `json:"raw_data"`
}

type RawData struct {
	Contract []Contract `json:"contract"`
}

type TransactionResponse struct {
	Data    []Transaction `json:"data"`
	Success bool          `json:"success"`
	Meta    Meta          `json:"meta"`
}

type Meta struct {
	At          int64  `json:"at"`
	Fingerprint string `json:"fingerprint"`
	Pagesize    int    `json:"page_size"`
}

type Contract struct {
	Parameter Parameter `json:"parameter"`
	Type      string    `json:"type"`
}

type Parameter struct {
	Value Value `json:"value"`
}

type Value struct {
	Amount       int64  `json:"amount"`
	OwnerAddress string `json:"owner_address"`
	ToAddress    string `json:"to_address"`
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID         string `json:"id"`
		Time       int64  `json:"time"`
		InAddress  string `json:"in_address"`
		OutAddress string `json:"out_address"`
		Amount     int64  `json:"amount"`
	}{
		ID:         t.ID,
		Time:       t.Timestamp,
		InAddress:  t.RawData.Contract[0].Parameter.Value.ToAddress,
		OutAddress: t.RawData.Contract[0].Parameter.Value.OwnerAddress,
		Amount:     t.RawData.Contract[0].Parameter.Value.Amount,
	})
}
