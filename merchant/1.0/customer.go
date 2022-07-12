package merchant

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/adless-tech/go-revolut/merchant/1.0/request"
)

type CustomerService struct {
	apiKey string
	domain string
}

type CreateCustomerReq struct {
	FullName     string `json:"full_name"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
}

type CreateCustomerResp struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Create creates a customer
// https://developer.revolut.com/docs/api-reference/merchant/#tag/Customers/operation/createCustomer
func (a *CustomerService) Create(req *CreateCustomerReq) (*CreateCustomerResp, error) {
	resp, statusCode, err := request.New(request.Config{
		Method:      http.MethodPost,
		Url:         fmt.Sprintf("%s/api/1.0/customers", a.domain),
		ApiKey:      a.apiKey,
		Body:        req,
		ContentType: request.ContentType_APPLICATION_JSON,
	})
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, errors.New(string(resp))
	}

	var r *CreateCustomerResp
	if err := json.Unmarshal(resp, r); err != nil {
		return nil, err
	}

	return r, nil
}
