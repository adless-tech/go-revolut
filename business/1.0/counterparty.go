package business

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/adless-tech/go-revolut/business/1.0/request"
)

type CounterpartyService struct {
	accessToken string
	sandbox     bool

	err error
}

type CounterpartyProfileType string

const (
	CounterpartyProfileType_BUSINESS CounterpartyProfileType = "business"
	CounterpartyProfileType_PERSONAL CounterpartyProfileType = "personal"
)

type CounterpartyType string

const (
	CounterpartyType_SELF     CounterpartyType = "self"
	CounterpartyType_REVOLUT  CounterpartyType = "revolut"
	CounterpartyType_EXTERNAL CounterpartyType = "external"
)

type RevolutCounterpartyReq struct {
	// the type of the Revolut profile, business or personal
	ProfileType CounterpartyProfileType `json:"profile_type,omitempty"`
	// an optional name. Provide only with personal profile_type
	Name string `json:"name,omitempty,omitempty"`
	// an optional phone number of the counterparty. Provide only with personal profile_type.
	Phone string `json:"phone,omitempty,omitempty"`
	// an optional email address of an admin of a public Revolut Business account. Provide only with business profile_type.
	Email string `json:"email,omitempty,omitempty"`
}

type NonRevolutCounterpartyReq struct {
	// an optional name of the external company counterparty, this field must exist when individual_name does not
	CompanyName  string                                `json:"company_name,omitempty,omitempty"`
	InvidualName NonRevolutCounterpartyReqInvidualName `json:"invidual_name,omitempty,omitempty"`
	// the country of the bank
	BankCountry string `json:"bank_country,omitempty"`
	// the currency of a counterparty's account
	Currency string `json:"currency,omitempty"`
	// bank account number
	AccountNo string `json:"account_no,omitempty"`
	// sort code
	SortCode string `json:"sort_code,omitempty"`
	// routing transit number
	RoutingNumber string `json:"routing_number,omitempty"`
	// an optional email address of the beneficiary
	Email string `json:"email,omitempty,omitempty"`
	// an optional phone number of the beneficiary
	Phone   string                           `json:"phone,omitempty,omitempty"`
	Address NonRevolutCounterpartyReqAddress `json:"address,omitempty"`
}

type NonRevolutCounterpartyReqInvidualName struct {
	// an optional first name of the external individual counterparty, this field must exist when company_name does not
	FirstName string `json:"first_name,omitempty,omitempty"`
	// an optional last name of the external individual counterparty, this field must exist when company_name does not
	LastName string `json:"last_name,omitempty,omitempty"`
}

type NonRevolutCounterpartyReqAddress struct {
	// an optional address line 1 of the counterparty
	StreetLine1 string `json:"street_line1,omitempty,omitempty"`
	// an optional address line 2 of the counterparty
	StreetLine2 string `json:"street_line2,omitempty,omitempty"`
	// an optional region of the counterparty
	Region string `json:"region,omitempty,omitempty"`
	// an optional postal code of the counterparty
	Postcode string `json:"postcode,omitempty,omitempty"`
	// an optional city of the counterparty
	City string `json:"city,omitempty,omitempty"`
	// an optional the bankCountry of the counterparty
	Country string `json:"country,omitempty,omitempty"`
}

type CounterpartyState string

const (
	CounterpartyState_ACTIVE   CounterpartyState = "created"
	CounterpartyState_INACTIVE CounterpartyState = "deleted"
)

type CounterpartyResp struct {
	// the ID of the counterparty
	Id string `json:"id,omitempty"`
	// the name of the counterparty
	Name string `json:"name,omitempty"`
	// the phone number of the counterparty
	Phone string `json:"phone,omitempty"`
	// the type of the Revolut profile, business or personal
	ProfileType CounterpartyProfileType `json:"profile_type,omitempty"`
	// the country of the bank
	Country string `json:"country,omitempty"`
	// the state of the counterparty, one of created, deleted
	State CounterpartyState `json:"state,omitempty"`
	// the instant when the counterparty was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// the instant when the counterparty was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// the list of public accounts of this counterparty
	Accounts []CounterpartyRespAccount `json:"accounts,omitempty"`
}

type CounterpartyRecipientCharges string

const (
	CounterpartyRecipientCharges_NO       CounterpartyRecipientCharges = "no"
	CounterpartyRecipientCharges_EXPECTED CounterpartyRecipientCharges = "expected"
)

type CounterpartyRespAccount struct {
	// the ID of a counterparty's account
	Id string `json:"id,omitempty"`
	// the currency of a counterparty's account
	Currency string `json:"currency,omitempty"`
	// the type of account, revolut or external
	Type string `json:"type,omitempty"`
	// bank account number
	AccountNo string `json:"account_no,omitempty"`
	// IBAN
	Iban string `json:"iban,omitempty"`
	// sort code
	SortCode    string `json:"sort_code,omitempty"`
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
	BankCountry string `json:"bank_country,omitempty"`
	// routing transit number
	RoutingNumber string `json:"routing_number,omitempty"`
	// BIC
	Bic string `json:"bic,omitempty"`
	// indicates the possibility of the recipient charges: no or expected
	RecipientCharges CounterpartyRecipientCharges `json:"recipient_charges,omitempty"`
}

// AddRevolut: You can create a counterparty for an existing Revolut user.
// doc: https://revolut-engineering.github.io/api-docs/#business-api-business-api-counterparties-add-revolut-counterparty
func (c *CounterpartyService) AddRevolut(revolutCounterparty *RevolutCounterpartyReq) (*CounterpartyResp, error) {
	if c.err != nil {
		return nil, c.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         "https://b2b.revolut.com/api/1.0/counterparty",
		AccessToken: c.accessToken,
		Sandbox:     c.sandbox,
		Body:        revolutCounterparty,
		ContentType: request.ContentType_APPLICATION_JSON,
	})
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &CounterpartyResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// AddNonRevolut: You can create a counterparty for an non-Revolut bank account.
// doc: https://revolut-engineering.github.io/api-docs/#business-api-business-api-counterparties-add-non-revolut-counterparty
func (c *CounterpartyService) AddNonRevolut(nonRevolutCounterparty *NonRevolutCounterpartyReq) (*CounterpartyResp, error) {
	if c.err != nil {
		return nil, c.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         "https://b2b.revolut.com/api/1.0/counterparty",
		AccessToken: c.accessToken,
		Sandbox:     c.sandbox,
		ContentType: request.ContentType_APPLICATION_JSON,
		Body:        nonRevolutCounterparty,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &CounterpartyResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Delete: This endpoint deletes a counterparty with the given ID.
// Once a counterparty is deleted no payments can be made to it.
// doc: https://revolut-engineering.github.io/api-docs/#business-api-business-api-counterparties-delete-counterparty
func (c *CounterpartyService) Delete(id string) error {
	if c.err != nil {
		return c.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodDelete,
		Url:         fmt.Sprintf("https://b2b.revolut.com/api/1.0/counterparty/%s", id),
		AccessToken: c.accessToken,
		Sandbox:     c.sandbox,
		Body:        nil,
	})

	if statusCode != http.StatusNoContent {
		return errors.New(string(resp))
	}

	if err != nil {
		return err
	}

	return nil
}

// WithId: This endpoint retrieves a counterparty by ID.
// doc https://revolut-engineering.github.io/api-docs/#business-api-business-api-counterparties-get-counterparty
func (c *CounterpartyService) WithId(id string) (*CounterpartyResp, error) {
	if c.err != nil {
		return nil, c.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         fmt.Sprintf("https://b2b.revolut.com/api/1.0/counterparty/%s", id),
		AccessToken: c.accessToken,
		Sandbox:     c.sandbox,
		Body:        nil,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &CounterpartyResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// List: This endpoint retrieves all your counterparties.
// doc: https://revolut-engineering.github.io/api-docs/#business-api-business-api-counterparties-get-counterparties
func (c *CounterpartyService) List() ([]*CounterpartyResp, error) {
	if c.err != nil {
		return nil, c.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         "https://b2b.revolut.com/api/1.0/counterparties",
		AccessToken: c.accessToken,
		Sandbox:     c.sandbox,
		Body:        nil,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := []*CounterpartyResp{}
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r, nil
}
