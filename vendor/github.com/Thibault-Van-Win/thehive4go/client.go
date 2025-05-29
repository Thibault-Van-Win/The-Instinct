package thehive4go

import (
	"crypto/tls"
	"net/http"

	"github.com/Thibault-Van-Win/thehive4go/alert"
	cases "github.com/Thibault-Van-Win/thehive4go/case"
	"github.com/Thibault-Van-Win/thehive4go/query"
)

type APIClient struct {
	http   *http.Client
	config Config

	// Core services
	query query.Service

	// Resource services
	Alerts alert.Service
	Cases  cases.Service
}

func NewAPIClient(config Config) *APIClient {
	httpClient := createHTTPClient(config)

	client := &APIClient{
		http:   httpClient,
		config: config,
	}

	client.query = *query.NewService(client)
	client.Alerts = *alert.NewService(client)
	client.Cases = *cases.NewService(client)

	return client
}

// Create HTTP client based on config
func createHTTPClient(config Config) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.SkipTLSVerification,
		},
	}

	return &http.Client{
		Transport: transport,
	}
}

// Decorate http requests
func (c *APIClient) DO(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	return c.http.Do(req)
}

func (c *APIClient) BaseUrl() string {
	return c.config.URL
}

func (c *APIClient) Query() *query.Service {
	return &c.query
}
