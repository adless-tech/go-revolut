package business

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/adless-tech/go-revolut/business/1.0/request"
)

type PaymentDraftService struct {
	accessToken string
	sandbox     bool

	err error
}

type PaymentDraftReq struct {
	// an optional title of payment
	Title string `json:"title,omitempty"`
	// an optional future date/time
	ScheduleFor string `json:"schedule_for,omitempty"`
	// a list of planned transactions
	Payments []PaymentDraftPayment `json:"payments,omitempty"`
}

type PaymentDraftResp struct {
	// the ID of the created draft payment
	Id string `json:"id,omitempty"`
}

type PaymentDraftPayment struct {
	// the transaction currency
	Currency string `json:"currency,omitempty"`
	// the transaction amount
	Amount int `json:"amount,omitempty"`
	// the ID of the account to pay from (must be the same for all payments json)
	AccountId string                      `json:"account_id,omitempty"`
	Receiver  PaymentDraftPaymentReceiver `json:"receiver,omitempty,omitempty"`
	// a mandatory textual reference shown on the transaction
	Reference string `json:"reference,omitempty"`
}

type PaymentDraftPaymentReceiver struct {
	// the ID of the receiving counterparty
	CounterpartyId string `json:"counterparty_id,omitempty"`
	// an optional ID of the receiving counterparty's account, can be own account (only for internal counterparties)
	AccountId string `json:"account_id,optional,omitempty"`
}

type PaymentDrafts struct {
	// a list of payments
	PaymentOrders []PaymentOrder `json:"payment_orders,omitempty"`
}

type PaymentOrder struct {
	// the ID of the draft payment
	Id string `json:"id,omitempty"`
	// an optional future date/time
	ScheduledFor string `json:"scheduled_for,optional,omitempty"`
	// an optional title of payment
	Title string `json:"title,optional,omitempty"`
	// count of payments in current draft
	PaymentsCount int `json:"payments_count,omitempty"`
}

type PaymentDraftDetail struct {
	// an optional future date/time
	ScheduledFor string `json:"scheduled_for,omitempty"`
	// an optional title of payment
	Title string `json:"title,omitempty"`
	// a list of payments
	Payments []PaymentDraftDetailPayment `json:"payments,omitempty"`
}

type PaymentDraftState string

const (
	PaymentDraftState_PENDING   PaymentDraftState = "PENDING"
	PaymentDraftState_COMPLETE  PaymentDraftState = "COMPLETED"
	PaymentDraftState_DECLINE   PaymentDraftState = "DECLINED"
	PaymentDraftState_FAILED    PaymentDraftState = "FAILED"
	PaymentDraftState_CREATED   PaymentDraftState = "CREATED"
	PaymentDraftState_REVERTED  PaymentDraftState = "REVERTED"
	PaymentDraftState_CANCELLED PaymentDraftState = "CANCELLED"
	PaymentDraftState_DELETED   PaymentDraftState = "DELETED"
)

type PaymentDraftDetailPayment struct {
	Id     string           `json:"id,omitempty"`
	Amount ExchangeRateResp `json:"amount,omitempty"`
	// the ID of the account to pay from
	AccountId string `json:"account_id,omitempty"`
	// an optional textual reference shown on the transaction
	Reference string                      `json:"reference,omitempty,omitempty"`
	Receiver  PaymentDraftPaymentReceiver `json:"receiver,omitempty"`
	// the state of the transaction, one of CREATED, PENDING, COMPLETED, REVERTED, DECLINED, CANCELLED, FAILED, DELETED
	State PaymentDraftState `json:"state,omitempty"`
	// an optional textual description of state reason
	Reason string `json:"reason,omitempty,omitempty"`
	// an optional textual description of error
	ErrorMessage string `json:"error_message,omitempty,omitempty"`
	// explanation of conversation process
	CurrentChargeOptions ExchangeRateResp `json:"current_charge_options,omitempty"`
}

// Create:
// doc: https://revolut-engineering.github.io/api-docs/business-api/#payment-drafts-create-a-payment-draft
func (e *PaymentDraftService) Create(paymentDraftReq *PaymentDraftReq) (*PaymentDraftResp, error) {
	if e.err != nil {
		return nil, e.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         "https://b2b.revolut.com/api/1.0/payment-drafts",
		AccessToken: e.accessToken,
		Sandbox:     e.sandbox,
		Body:        paymentDraftReq,
		ContentType: request.ContentType_APPLICATION_JSON,
	})
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &PaymentDraftResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// List:
// doc: https://revolut-engineering.github.io/api-docs/business-api/#get-payment-drafts
func (e *PaymentDraftService) List() (*PaymentDrafts, error) {
	if e.err != nil {
		return nil, e.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         "https://b2b.revolut.com/api/1.0/payment-drafts",
		AccessToken: e.accessToken,
		Sandbox:     e.sandbox,
	})
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &PaymentDrafts{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// WithId:
// doc: https://revolut-engineering.github.io/api-docs/business-api/#get-payment-drafts-get-payment-draft-by-id
func (e *PaymentDraftService) WithId(id string) (*PaymentDraftDetailPayment, error) {
	if e.err != nil {
		return nil, e.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodGet,
		Url:         fmt.Sprintf("https://b2b.revolut.com/api/1.0/payment-drafts/%s", id),
		AccessToken: e.accessToken,
		Sandbox:     e.sandbox,
	})
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	r := &PaymentDraftDetailPayment{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Delete:
// doc: https://revolut-engineering.github.io/api-docs/business-api/#get-payment-drafts-delete-payment-draft
func (e *PaymentDraftService) Delete(id string) error {
	if e.err != nil {
		return e.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodDelete,
		Url:         fmt.Sprintf("https://b2b.revolut.com/api/1.0/payment-drafts/%s", id),
		AccessToken: e.accessToken,
		Sandbox:     e.sandbox,
	})
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return errors.New(string(resp))
	}

	return nil
}
