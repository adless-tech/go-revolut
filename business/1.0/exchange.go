package business

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/adless-tech/go-revolut/business/1.0/request"
)

type ExchangeService struct {
	accessToken string
	sandbox     bool

	err error
}

type ExchangeRateReq struct {
	// the currency you would like to exchange from
	From string
	// the currency you would like to exchange to
	To string
	// exchange amount, default is 1.00
	Amount float64
}

type ExchangeRateResp struct {
	// information about the currency to exchange from
	From Amount `json:"from,omitempty"`
	// information about the currency to exchange to
	To Amount `json:"to,omitempty"`
	// exchange rate
	Rate float64 `json:"rate,omitempty"`
	// fee for the operation
	Fee Amount `json:"fee,omitempty"`
	// date of proposed exchange rate
	RateDate time.Time `json:"rate_date,omitempty"`
}

type Amount struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

type ExchangeReq struct {
	// information about the account you want to exchange from
	From ExchangeAmount `json:"from,omitempty"`
	// information about the account you want to exchange to
	To ExchangeAmount `json:"to,omitempty"`
	// a user-provided exchange reference
	Reference string `json:"reference,omitempty"`
	// a unique value used to handle duplicates submitted as a result of lost connection or another client error (40 characters max)
	RequestId string `json:"request_id,omitempty"`
}
type ExchangeAmount struct {
	// the account ID
	AccountId string  `json:"account_id,omitempty"`
	Amount    float64 `json:"amount,omitempty,omitempty"`
	Currency  string  `json:"currency,omitempty"`
}

type ExchangeResp struct {
	// the ID of transaction
	Id string `json:"id,omitempty"`
	// is always exchange
	State string `json:"state,omitempty"`
	// reason code for declined or failed transaction state
	ReasonCode string `json:"reason_code,omitempty"`
	// the instant when the transaction was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// the instant when the transaction was completed
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// Rate:
// doc: https://revolut-engineering.github.io/api-docs/business-api/#exchanges-get-exchange-rates
func (e *ExchangeService) Rate(exchangeRateReq *ExchangeRateReq) (*ExchangeRateResp, error) {
	if e.err != nil {
		return nil, e.err
	}

	params := url.Values{}
	params.Add("from", exchangeRateReq.From)
	params.Add("to", exchangeRateReq.To)
	params.Add("amount", fmt.Sprintf("%0.2f", exchangeRateReq.Amount))

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         fmt.Sprintf("https://b2b.revolut.com/api/1.0/rate?%s", params.Encode()),
		AccessToken: e.accessToken,
		Sandbox:     e.sandbox,
	})
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &ExchangeRateResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Exchange: To check the exchange rate and fees for the operation, please use the /rate endpoint.
// doc: https://revolut-engineering.github.io/api-docs/business-api/#exchanges-exchange-currency
func (e *ExchangeService) Exchange(exchangeReq *ExchangeReq) (*ExchangeResp, error) {
	if e.err != nil {
		return nil, e.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         "https://b2b.revolut.com/api/1.0/exchange",
		AccessToken: e.accessToken,
		Sandbox:     e.sandbox,
		Body:        exchangeReq,
		ContentType: request.ContentType_APPLICATION_JSON,
	})
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &ExchangeResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}
