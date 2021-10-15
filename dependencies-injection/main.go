package main

func main() {

	// 0. Using Multiple Initialization Functions
	/*

		func NewAPIRequest() (*http.Request, error) {
			return http.NewRequest(http.MethodGet, ServerURI, nil)
		}

		func NewRequestWithAuth(key string) (*http.Request, error) {
			req, err := NewAPIRequest()
			if err != nil {
				return nil, err
			}
			req.Header.Add("authorization", key)

			return req, err
		}

	*/

	// 1. pass dependency via config
	_ = NewClientWithConfig(ClientConfig{})

	// 2. Using Variadic Functions with Option Functions
	_ = NewClient(WithURL("https://secure.example"))

	// 3. pass interfaces, return structs
	storage := Storage{}
	_ = NewCustomer(storage)
}

type Storage struct {
}

func (s Storage) Data() []string {
	return []string{"Dummy"}
}

type ClientConfig struct {
	Url string
}
type Client struct {
	url string
}

func NewClientWithConfig(cfg ClientConfig) *Client {
	return &Client{url: cfg.Url}
}

func NewClient(options ...func(client *Client)) *Client {
	client := &Client{}
	for _, opt := range options {
		opt(client)
	}

	return client
}

func WithURL(url string) func(client *Client) {
	return func(client *Client) {
		client.url = url
	}
}

type Share interface {
	Data() []string
}
type Account struct {
	Info []string
}

func NewCustomer(share Share) *Account {
	return &Account{Info: share.Data()}
}
