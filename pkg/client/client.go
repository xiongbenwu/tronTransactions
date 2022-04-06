package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
	"transactions/pkg/tron"
)

const updateRate = 5

type TRONClient struct {
	client       *http.Client
	mu           sync.RWMutex
	mainAddress  string
	apiURL       *url.URL
	transactions []tron.Transaction
	lastUpdate   int64
}

func (c *TRONClient) GetTransactions(address string) ([]tron.Transaction, error) {
	c.mu.RLock()
	transactions := c.transactions
	c.mu.RUnlock()

	var matched []tron.Transaction

	for _, t := range transactions {
		if t.RawData.Contract.Parameter.Value.OwnerAddress == address ||
			t.RawData.Contract.Parameter.Value.ToAddress == address {
			matched = append(matched, t)
		}
	}

	return matched, nil
}

func NewClient(address string, apiURL string) (*TRONClient, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	c := &TRONClient{
		client:      &http.Client{},
		mu:          sync.RWMutex{},
		mainAddress: address,
		apiURL:      u,
	}

	return c, nil
}

// Run receives all transactions at start and updates array every 5 seconds.
func (c *TRONClient) Run() error {
	c.mu.Lock()

	transactions, err := c.getTransactionsRecursive("")
	if err != nil {
		return err
	}

	c.transactions = transactions
	c.mu.Unlock()

	ticker := time.NewTicker(updateRate * time.Second)

	for {
		<-ticker.C

		transactions, err := c.getTransactionsRecursive("")
		if err != nil {
			ticker.Stop()

			return err
		}

		c.mu.Lock()
		c.transactions = append(c.transactions, transactions...)
		c.mu.Unlock()
	}
}

func (c *TRONClient) getTransactionsRecursive(fingerprint string) ([]tron.Transaction, error) {
	req, err := c.newGetTransactionRequest(fingerprint, fmt.Sprintf("%d", c.lastUpdate))
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var trResponse tron.TransactionResponse

	err = json.Unmarshal(body, &trResponse)
	if err != nil {
		return nil, err
	}

	var transferTransactions []tron.Transaction

	for _, t := range trResponse.Data {
		if t.RawData.Contract.Type == "TransferContract" {
			transferTransactions = append(transferTransactions, t)
		}
	}

	if trResponse.Meta.Fingerprint != "" {
		t, err := c.getTransactionsRecursive(trResponse.Meta.Fingerprint)
		if err != nil {
			return nil, err
		}
		transferTransactions = append(transferTransactions, t...)
	} else {
		c.lastUpdate = trResponse.Meta.At
	}

	return transferTransactions, nil
}

func (c *TRONClient) newGetTransactionRequest(fingerprint, minTimestamp string) (*http.Request, error) {
	order := "block_timestamp,asc"
	u := *c.apiURL
	q := u.Query()
	q.Set("limit", "200")

	if fingerprint != "" {
		q.Add("fingerprint", fingerprint)
	}

	q.Add("order_by", order)
	q.Add("min_timestamp", minTimestamp)

	u.RawQuery = q.Encode()
	u.Path = fmt.Sprintf("%saccounts/%s/transactions", u.Path, c.mainAddress)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")

	return req, nil
}
