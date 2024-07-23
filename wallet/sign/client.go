package sign

import (
	"fmt"
	"github.com/pkg/errors"

	gresty "github.com/go-resty/resty/v2"
)

var errSolSignHTTPError = errors.New("solana chain http error")

type SolSignClient interface {
	GenerateAddress(uint64) (*AccountInfoRep, error)
	PrepareAccount(*PrepareAccountReq) (*PrepareAccountRep, error)
	SignTransaction(*TransactionReq) (*TransactionRep, error)
}

type Client struct {
	client *gresty.Client
}

func NewSolSignClient(url string) (*Client, error) {
	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errSolSignHTTPError)
		}
		return nil
	})
	return &Client{
		client: client,
	}, nil
}

func (c *Client) GenerateAddress(addressNum uint64) (accountInfoRep *AccountInfoRep, err error) {
	var accountInfoRetRep AccountInfoRep
	_, err = c.client.R().
		SetBody(map[string]interface{}{"address_num": addressNum}).
		SetResult(&accountInfoRetRep).
		Post("/generateAddress")
	if err != nil {
		return nil, fmt.Errorf("genearate address fail: %w", err)
	}
	return &accountInfoRetRep, nil
}

func (c *Client) PrepareAccount(prepareAccountReq *PrepareAccountReq) (prepareAccountRep *PrepareAccountRep, err error) {
	var prepareAccountReponse PrepareAccountRep
	_, err = c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(prepareAccountReq).
		SetResult(&prepareAccountReponse).
		Post("/prepareAccount")
	if err != nil {
		return nil, fmt.Errorf("prepare account fail: %w", err)
	}
	return &prepareAccountReponse, nil
}

func (c *Client) SignTransaction(transactionReq *TransactionReq) (transactionRep *TransactionRep, err error) {
	var transactionReponse TransactionRep
	_, err = c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(transactionReq).
		SetResult(&transactionReponse).
		Post("/signTransaction")
	if err != nil {
		return nil, fmt.Errorf("prepare account fail: %w", err)
	}
	return &transactionReponse, nil
}
