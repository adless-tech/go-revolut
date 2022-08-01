package business

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adless-tech/go-revolut/business/1.0/request"
)

type AccountService struct {
	accessToken string
	sandbox     bool

	err error
}

type AccountState string

const (
	AccountState_ACTIVE   AccountState = "active"
	AccountState_INACTIVE AccountState = "inactive"
)

type AccountResp struct {
	// the account ID
	Id string `json:"id,omitempty"`
	// the account name
	Name string `json:"name,omitempty"`
	// the available balance
	Balance float64 `json:"balance,omitempty"`
	// the account currency
	Currency string `json:"currency,omitempty"`
	// the account state, one of active, inactive
	State AccountState `json:"state,omitempty"`
	// determines if the account is visible to other businesses on Revolut
	Public bool `json:"public,omitempty"`
	// the instant when the account was created
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// the instant when the account was last updated
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type AccountSchema string

const (
	AccountSchema_CHAPS           AccountSchema = "chaps"
	AccountSchema_BACS            AccountSchema = "bacs"
	AccountSchema_FASTER_PAYMENTS AccountSchema = "faster_payments"
	AccountSchema_SEPA            AccountSchema = "sepa"
	AccountSchema_SWIFT           AccountSchema = "swift"
	AccountSchema_ACH             AccountSchema = "ach"
)

type AccountDetailResp struct {
	// IBAN
	Iban string `json:"iban,omitempty"`
	// BIC
	Bic string `json:"bic,omitempty"`
	// the account number
	AccountNo string `json:"account_no,omitempty"`
	// the sort code
	SortCode string `json:"sort_code,omitempty"`
	// the routing number
	RoutingNumber string `json:"routing_number,omitempty"`
	// the beneficiary name
	Beneficiary        string             `json:"beneficiary,omitempty"`
	BeneficiaryAddress BeneficiaryAddress `json:"beneficiary_address,omitempty"`
	// the country of the bank
	BankCountry string `json:"bank_country,omitempty"`
	// determines if this account address is pooled or unique
	Pooled bool `json:"pooled,omitempty"`
	// the reference of the pooled account
	UniqueReference string `json:"unique_reference,omitempty"`
	// the list of supported schemes, possible values: chaps, bacs, faster_payments, sepa, swift, ach
	Schemes       []AccountSchema `json:"schemes,omitempty"`
	EstimatedTime EstimatedTime   `json:"estimated_time,omitempty"`
}

type AccountUnit string

const (
	AccountUnit_DAYS  AccountUnit = "days"
	AccountUnit_HOURS AccountUnit = "hours"
)

type EstimatedTime struct {
	// the unit of the inbound transfer time estimate, possible values: days, hours
	Unit AccountUnit `json:"unit,omitempty"`
	// the maximum estimate
	Min int `json:"min,omitempty"`
	// the minimum estimate
	Max int `json:"max,omitempty"`
}

type BeneficiaryAddress struct {
	// the address line 1 of the beneficiary
	StreetLine1 string `json:"street_line1,omitempty"`
	// the address line 2 of the beneficiary
	StreetLine2 string `json:"street_line2,omitempty"`
	// the region of the beneficiary
	Region string `json:"region,omitempty"`
	// the city of the beneficiary
	City string `json:"city,omitempty"`
	// the country of the beneficiary
	Country string `json:"country,omitempty"`
	// the postal code of the beneficiary
	Postcode string `json:"postcode,omitempty"`
}

// List: This endpoint retrieves your accounts.
// doc: https://revolut-engineering.github.io/api-docs/#business-api-business-api-accounts-get-accounts
func (a *AccountService) List() ([]*AccountResp, error) {
	if a.err != nil {
		return nil, a.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         "https://b2b.revolut.com/api/1.0/accounts",
		AccessToken: a.accessToken,
		Sandbox:     a.sandbox,
		Body:        nil,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	var r []*AccountResp
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r, nil
}

// WithId: This endpoint retrieves one of your accounts by ID.
// doc: https://revolut-engineering.github.io/api-docs/#business-api-business-api-accounts-get-account
func (a *AccountService) WithId(id string) (*AccountResp, error) {
	if a.err != nil {
		return nil, a.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         fmt.Sprintf("https://b2b.revolut.com/api/1.0/accounts/%s", id),
		AccessToken: a.accessToken,
		Sandbox:     a.sandbox,
		Body:        nil,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &AccountResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// DetailWithId: This endpoint retrieves individual account details.
// doc: https://revolut-engineering.github.io/api-docs/business-api/#accounts-get-account-details
func (a *AccountService) DetailWithId(id string) ([]*AccountDetailResp, error) {
	if a.err != nil {
		return nil, a.err
	}
	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         fmt.Sprintf("https://b2b.revolut.com/api/1.0/accounts/%s/bank-details", id),
		AccessToken: a.accessToken,
		Sandbox:     a.sandbox,
		Body:        nil,
	})
	if err != nil {
		return []*AccountDetailResp{}, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := []*AccountDetailResp{}
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r, nil
}
