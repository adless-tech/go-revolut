package business

import (
	"errors"
	"net/http"
	"time"

	"github.com/adless-tech/go-revolut/business/1.0/request"
)

type WebhookService struct {
	accessToken string
	sandbox     bool

	err error
}

type TransactionStateChangedEvent struct {
	// the event name
	Event string `json:"event,omitempty"`
	// the event time
	Timestamp time.Time                        `json:"timestamp,omitempty"`
	Data      TransactionStateChangedEventData `json:"data,omitempty"`
}

type TransactionStateChangedEventData struct {
	// the ID of the transaction
	ID string `json:"id,omitempty"`
	// previous state of the transaction
	OldState string `json:"old_state,omitempty"`
	// new state of the transaction
	NewState string `json:"new_state,omitempty"`
}

type TransactionCreatedEvent struct {
	// the event name
	Event string `json:"event,omitempty"`
	// the event time
	Timestamp time.Time                   `json:"timestamp,omitempty"`
	Data      TransactionCreatedEventData `json:"data,omitempty"`
}

type TransactionCreatedEventData struct {
	// the ID of transaction
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	// the client provided request ID
	RequestId string `json:"request_id,omitempty"`
	// the transction state: pending, completed, declined or failed
	State PaymentState `json:"state,omitempty"`
	// an optional reason code for declined or failed transaction state
	ReasonCode string `json:"reason_code,omitempty"`
	// the instant when the transaction was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// the instant when the transaction was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// the instant when the transaction was completed, mandatory for completed state only
	CompletedAt time.Time `json:"completed_at,omitempty"`
	// an optional date when the transaction was scheduled for
	ScheduledFor string `json:"scheduled_for,omitempty"`
	// a user provided payment reference
	Reference string `json:"reference,omitempty"`
	// the legs of transaction, there'll be 2 legs between your Revolut accounts and 1 leg in other cases
	Legs []TransactionLeg `json:"legs,omitempty"`
}

// Set:
// doc: https://revolut-engineering.github.io/api-docs/business-api/#web-hooks-setting-up-a-web-hook
func (p *WebhookService) Set(url string) error {
	if p.err != nil {
		return p.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         "https://b2b.revolut.com/api/1.0/webhook",
		AccessToken: p.accessToken,
		Sandbox:     p.sandbox,
		Body: struct {
			// call back endpoint of the client system, https is the supported protocol
			Url string `json:"url,omitempty"`
		}{Url: url},
		ContentType: request.ContentType_APPLICATION_JSON,
	})
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return errors.New(string(resp))
	}

	return nil
}

// Delete: Use this API request to delete a web-hook
// doc: https://revolut-engineering.github.io/api-docs/business-api/#web-hooks-setting-up-a-web-hook
func (p *WebhookService) Delete() error {
	if p.err != nil {
		return p.err
	}

	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodDelete,
		Url:         "https://b2b.revolut.com/api/1.0/webhook",
		AccessToken: p.accessToken,
		Sandbox:     p.sandbox,
	})
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return errors.New(string(resp))
	}

	return nil
}
