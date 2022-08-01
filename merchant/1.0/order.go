package merchant

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/adless-tech/go-revolut/merchant/1.0/request"
)

type OrderService struct {
	apiKey string
	domain string
}

type OrderType string

const (
	OrderType_PAYMENT    OrderType = "PAYMENT"
	OrderType_REFUND     OrderType = "REFUND"
	OrderType_CHARGEBACK OrderType = "CHARGEBACK"
)

type Funding string

const (
	Funding_CREDIT  Funding = "CREDIT"
	Funding_DEBIT   Funding = "DEBIT"
	Funding_PREPAID Funding = "PREPAID"
)

type TreeDsState string

const (
	TreeDSState_VERIFIED  TreeDsState = "VERIFIED"
	TreeDSState_FAILED    TreeDsState = "FAILED"
	TreeDSState_CHALLENGE TreeDsState = "CHALLENGE"
)

type CardType string

const (
	CardType_VISA       CardType = "VISA"
	CardType_MASTERCARD CardType = "MASTERCARD"
)

type RiskLevel string

const (
	RiskLevel_LOW  RiskLevel = "LOW"
	RiskLevel_HIGH RiskLevel = "HIGH"
)

type CvvVerification string

const (
	CvvVerification_MATCH         CvvVerification = "MATCH"
	CvvVerification_NOT_MATCH     CvvVerification = "NOT_MATCH"
	CvvVerification_INCORRECT     CvvVerification = "INCORRECT"
	CvvVerification_NOT_PROCESSED CvvVerification = "NOT_PROCESSED"
)

type CheckResult string

const (
	CheckResult_MATCH     CheckResult = "MATCH"
	CheckResult_NOT_MATCH CheckResult = "NOT_MATCH"
	CheckResult_N_A       CheckResult = "N_A"
	CheckResult_INVALID   CheckResult = "INVALID"
)

type OrderState string

const (
	OrderState_PENDING    OrderState = "PENDING"
	OrderState_PROCESSING OrderState = "PROCESSING"
	OrderState_AUTHORISED OrderState = "AUTHORISED"
	OrderState_COMPLETED  OrderState = "COMPLETED"
	OrderState_FAILED     OrderState = "FAILED"
)

type Amount struct {
	Value    int    `json:"value"`
	Currency string `json:"currency"`
}

type FeeType string

const (
	FeeType_FX        FeeType = "FX"
	FeeType_ACQUIRING FeeType = "ACQUIRING"
)

type Fee struct {
	// Fee amount
	Value int `json:"value,omitempty"`
	// Fee currency
	Currency string `json:"currency,omitempty"`
	// Fee type
	Type FeeType `json:"type,omitempty"`
}

type Payment struct {
	Type          string `json:"type,omitempty"`
	Amount        Amount `json:"amount,omitempty"`
	State         string `json:"state,omitempty"`
	CreatedDate   int64  `json:"created_date,omitempty"`
	UpdatedDate   int64  `json:"updated_date,omitempty"`
	CompletedDate int    `json:"completed_date,omitempty"`
	Card          Card   `json:"card,omitempty"`
}

type ThreeDs struct {
	// 3DS check result
	State TreeDsState `json:"state,omitempty"`
	// 3DS version
	Version int `json:"version,omitempty"`
}

type Check struct {
	// Confirms whether a proxy number is used or not
	Proxy bool `json:"proxy,omitempty"`
	// Confirms whether a VPN connection is used or not
	Vpn bool `json:"vpn,omitempty"`
	// Country name associated with the IP address of the card used
	CountryByIp string  `json:"country_by_ip,omitempty"`
	ThreeDs     ThreeDs `json:"three_ds,omitempty"`
	// Authorization code returned by the Processor
	AuthorizationCode string `json:"authorization_code,omitempty"`
	// CVV verification
	CvvVerification CvvVerification `json:"cvv_verification,omitempty"`
	// Address verification
	Address CheckResult `json:"address,omitempty"`
	// Postal code verification
	PostalCode CheckResult `json:"postal_code,omitempty"`
	// Cardholder verification
	CardHolder CheckResult `json:"card_holder,omitempty"`
}

type Card struct {
	// Card type
	CardType CardType `json:"card_type,omitempty"`
	// Card funding
	Funding Funding `json:"funding,omitempty"`
	// Card BIN
	CardBin string `json:"card_bin,omitempty"`
	// Card last four digits
	CardLastFour string `json:"card_last_four,omitempty"`
	// Card expiry date in the format of MM/YY
	CardExpiry string `json:"card_expiry,omitempty"`
	// Cardholder name
	CardholderName string         `json:"cardholder_name,omitempty"`
	Checks         Check          `json:"checks,omitempty"`
	RiskLevel      RiskLevel      `json:"risk_level,omitempty"`
	BillingAddress BillingAddress `json:"billing_address,omitempty"`
}

type BillingAddress struct {
	// Street line 1 information
	StreetLine1 string `json:"street_line_1,omitempty"`
	// Street line 2 information
	StreetLine2 string `json:"street_line_2,omitempty"`
	// Region name
	Region string `json:"region,omitempty"`
	// City name
	City string `json:"city,omitempty"`
	// Country associated with the address
	CountryCode string `json:"country_code,omitempty"`
	// Postcode associated with the address
	Postcode string `json:"postcode,omitempty"`
}

type OrderResp struct {
	// Order ID for a merchant
	Id string `json:"id,omitempty"`
	// Temporary ID for a customer
	PublicId string `json:"public_id,omitempty"`
	// Temporary ID for a customer
	Type OrderType `json:"type,omitempty"`
	// Order state
	State OrderState `json:"state,omitempty"`
	// Order creation date, measured in ms since the Unix epoch (UTC)
	CreatedDate int64 `json:"created_date,omitempty"`
	// Last update date, measured in ms since the Unix epoch (UTC)
	UpdatedDate int64 `json:"updated_date,omitempty"`
	// Order completion date, measured in ms since the Unix epoch (UTC)
	CompletedDate int64  `json:"completed_date,omitempty"`
	OrderAmount   Amount `json:"order_amount,omitempty"`
	// Merchant order ID
	MerchantOrderExtRef string `json:"merchant_order_ext_ref,omitempty"`
	// Merchant customer ID
	MerchantCustomerExtRef string `json:"merchant_customer_ext_ref,omitempty"`
	// Customer e-mail
	Email           string           `json:"email,omitempty"`
	SettledAmount   Amount           `json:"settled_amount,omitempty"`
	RefundedAmount  Amount           `json:"refunded_amount,omitempty"`
	Fees            []Fee            `json:"fees,omitempty"`
	Payments        []Payment        `json:"payments,omitempty"`
	Attempts        []AttemptRelated `json:"attempts,omitempty"`
	Related         []AttemptRelated `json:"related,omitempty"`
	ShippingAddress ShippingAddress  `json:"shipping_address,omitempty"`
	Phone           string           `json:"phone,omitempty"`
	CustomerID      string           `json:"customer_id,omitempty"`
}

type AttemptRelated struct {
	Id     string    `json:"id,omitempty"`
	Type   OrderType `json:"type,omitempty"`
	Amount Amount    `json:"amount,omitempty"`
}

type ShippingAddress struct {
	// Shipping address: Street line 1 information
	StreetLine1 string `json:"street_line_1,omitempty"`
	// Shipping address: Street line 2 information
	StreetLine2 string `json:"street_line_2,omitempty"`
	// Shipping address: Region name
	Region string `json:"region,omitempty"`
	// Shipping address: City name
	City string `json:"city,omitempty"`
	// Shipping address: Country associated with the address
	CountryCode string `json:"country_code,omitempty"`
	// Shipping address: Postcode associated with the address
	Postcode string `json:"postcode,omitempty"`
}

type CaptureMode string

const (
	CaptureMode_MANUAL    CaptureMode = "MANUAL"
	CaptureMode_AUTOMATIC CaptureMode = "AUTOMATIC"
)

type OrderReq struct {
	// Minor amount
	Amount int `json:"amount"`
	// Capture mode. If it is equal to null then AUTOMATIC is used
	CaptureMode CaptureMode `json:"capture_mode,omitempty"`
	// Merchant order ID
	MerchantOrderID string `json:"merchant_order_id,omitempty"`
	// Customer e-mail
	CustomerEmail string `json:"customer_email,omitempty"`
	// Order description
	Description string `json:"description,omitempty"`
	// Currency code
	Currency string `json:"currency,omitempty"`
	// Settlement currency. If it is equal to null then the payment is settled in transaction currency.
	SettlementCurrency string `json:"settlement_currency,omitempty"`
	// Merchant customer ID
	MerchantCustomerID string `json:"merchant_customer_id,omitempty"`
}

type RefundReq struct {
	// Minor amount
	Amount int `json:"amount,omitempty"`
	// Merchant order ID
	MerchantOrderID string `json:"merchant_order_id,omitempty"`
	// Order description
	Description string `json:"description,omitempty"`
	// Currency code
	Currency string `json:"currency,omitempty"`
}

type RefundResp struct {
	// Order ID for a merchant
	Id string `json:"id,omitempty"`
	// Order type
	Type OrderType `json:"type,omitempty"`
	// Order state
	State OrderState `json:"state,omitempty"`
	// Order creation date, measured in ms since the Unix epoch (UTC)
	CreatedDate int64 `json:"created_date,omitempty"`
	// Last update date, measured in ms since the Unix epoch (UTC)
	UpdatedDate int64 `json:"updated_date,omitempty"`
	// Order completion date, measured in ms since the Unix epoch (UTC)
	CompletedDate int64  `json:"completed_date,omitempty"`
	OrderAmount   Amount `json:"order_amount,omitempty"`
	// Merchant customer ID
	MerchantCustomerExtRef string `json:"merchant_customer_ext_ref,omitempty"`
	// Customer e-mail
	Email   string           `json:"email,omitempty"`
	Related []AttemptRelated `json:"related,omitempty"`
}

// Create:
// doc: https://revolut-engineering.github.io/api-docs/merchant-api/#backend-api-backend-api-order-object-create-payment-order
func (a *OrderService) Create(orderReq *OrderReq) (*OrderResp, error) {
	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         fmt.Sprintf("%s/api/1.0/orders", a.domain),
		ApiKey:      a.apiKey,
		Body:        orderReq,
		ContentType: request.ContentType_APPLICATION_JSON,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("%s [status: %d]", string(resp), statusCode))
	}

	r := &OrderResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// WithId: If you would like to get information about the created order, please use the following request.
// doc: https://revolut-engineering.github.io/api-docs/merchant-api/#backend-api-backend-api-order-object-retrieve-order
func (a *OrderService) WithId(id string) (*OrderResp, error) {
	resp, statusCode, err := request.New(request.Config{
		Method: http.MethodGet,
		Url:    fmt.Sprintf("%s/api/1.0/orders/%s", a.domain, id),
		ApiKey: a.apiKey,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status: %d\n\nResponse: %s", statusCode, string(resp)))
	}

	r := &OrderResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Capture: Once the payment is authorised, the merchant needs to
// capture it in order for it to be sent into the processing stage.
// doc: https://revolut-engineering.github.io/api-docs/merchant-api/#backend-api-backend-api-order-object-capture-order
func (a *OrderService) Capture(id string) (*OrderResp, error) {
	resp, statusCode, err := request.New(request.Config{
		Method: http.MethodPost,
		Url:    fmt.Sprintf("%s/api/1.0/orders/%s/capture", a.domain, id),
		ApiKey: a.apiKey,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status: %d\n\nResponse: %s", statusCode, string(resp)))
	}

	r := &OrderResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Cancel: In case the payment has not been captured yet and the merchant decides
// to not proceed with the order, the order can be cancelled manually.
// doc: https://revolut-engineering.github.io/api-docs/merchant-api/#backend-api-backend-api-order-object-cancel-order
func (a *OrderService) Cancel(id string) (*OrderResp, error) {
	resp, statusCode, err := request.New(request.Config{
		Method: http.MethodPost,
		Url:    fmt.Sprintf("%s/api/1.0/orders/%s/cancel", a.domain, id),
		ApiKey: a.apiKey,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status: %d\n\nResponse: %s", statusCode, string(resp)))
	}

	r := &OrderResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}

// Refund: In case the customer requires a refund for a payment that has been already captured,
// the merchant can always issue a full or partial refund for a particular payment.
// doc: https://revolut-engineering.github.io/api-docs/merchant-api/#backend-api-backend-api-order-object-refund-order
func (a *OrderService) Refund(id string, refundReq *RefundReq) (*RefundResp, error) {
	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         fmt.Sprintf("%s/api/1.0/orders/%s/refund", a.domain, id),
		ApiKey:      a.apiKey,
		Body:        refundReq,
		ContentType: request.ContentType_APPLICATION_JSON,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status: %d\n\nResponse: %s", statusCode, string(resp)))
	}

	r := &RefundResp{}
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}
