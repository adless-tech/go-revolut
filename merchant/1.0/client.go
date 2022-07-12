package merchant

type Client struct {
	apiKey string
	domain string

	orderService    *OrderService
	webhookService  *WebhookService
	customerService *CustomerService
}

func newClient(apiKey string, domain string) *Client {
	return &Client{
		apiKey: apiKey,
		domain: apiKey,
		orderService: &OrderService{
			apiKey: apiKey,
			domain: domain,
		},
		webhookService: &WebhookService{
			apiKey: apiKey,
			domain: domain,
		},
		customerService: &CustomerService{
			apiKey: apiKey,
			domain: domain,
		},
	}
}

func NewProductionClient(apiKey string) *Client {
	return newClient(apiKey, "https://merchant.revolut.com")
}

func NewSandboxClient(apiKey string) *Client {
	return newClient(apiKey, "https://sandbox-merchant.revolut.com")
}

func (m *Client) Order() *OrderService {
	return m.orderService
}

func (m *Client) Webhook() *WebhookService {
	return m.webhookService
}

func (m *Client) Customer() *CustomerService {
	return m.customerService
}
