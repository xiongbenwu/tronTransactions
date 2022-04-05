package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"transactions/pkg/tron"
)

type TRONClient struct {
	client       *http.Client
	mu           sync.RWMutex
	mainAddress  string
	apiURL       *url.URL
	transactions []tron.Transaction
	lastUpdate   int64
}

func (c *TRONClient) GetTransactions() ([]tron.Transaction, error) {
	return c.transactions, nil
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

func (c *TRONClient) Run() error {
	err := c.getInitialTransactions()
	if err != nil {
		return err
	}
	go func() {
		c.scan()
	}()

	return nil
}

func (c *TRONClient) scan() error {

	return nil
}

func (c *TRONClient) getInitialTransactions() error {
	req, err := c.newGetTransactionRequest("")
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var trResponse tron.TransactionResponse

	err = json.Unmarshal(body, &trResponse)
	if err != nil {
		return err
	}

	if trResponse.Meta.Fingerprint != "" {
		err := c.getTransactionsRecursive(trResponse.Meta.Fingerprint, trResponse.Data)
		if err != nil {
			return err
		}
	}

	for _, t := range trResponse.Data {
		if t.RawData.Contract[0].Type == "TransferContract" {
			c.transactions = append(c.transactions, t)
		}
	}

	return nil
}

func (c *TRONClient) getTransactionsRecursive(fingerprint string, transactions []tron.Transaction) error {
	req, err := c.newGetTransactionRequest(fingerprint)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var trResponse tron.TransactionResponse

	err = json.Unmarshal(body, &trResponse)
	if err != nil {
		return err
	}

	transactions = append(transactions, trResponse.Data...)
	if trResponse.Meta.Fingerprint != "" {
		return c.getTransactionsRecursive(trResponse.Meta.Fingerprint, transactions)
	}

	return nil
}

func (c *TRONClient) newGetTransactionRequest(fingerprint string) (*http.Request, error) {
	u := *c.apiURL
	if fingerprint != "" {
		u.RawQuery = fmt.Sprintf("fingerprint=%s", fingerprint)
	}
	u.Path = fmt.Sprintf("%saccounts/%s/transactions", u.Path, c.mainAddress)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	return req, nil
}
