package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ClientInterface defines required methods from the main client
type BaseClient interface {
	DO(*http.Request) (*http.Response, error)
	BaseUrl() string
}

// Service handles query operations
type Service struct {
	client BaseClient
}

// NewService creates a new query service
func NewService(client BaseClient) *Service {
	return &Service{client: client}
}

func (s *Service) Execute(options ...Option) ([]byte, error) {

	opts := initOptions()

	for _, option := range options {
		option(&opts)
	}

	query, err := opts.buildQuery()
	if err != nil {
		return nil, fmt.Errorf("error building query: %v", err)
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %v", err)
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, s.client.BaseUrl()+"/api/v1/query", bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Execute request
	resp, err := s.client.DO(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad response status: %s, body: %s", resp.Status, string(body))
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
