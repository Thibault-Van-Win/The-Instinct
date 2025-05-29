package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Thibault-Van-Win/thehive4go/query"
	"github.com/Thibault-Van-Win/thehive4go/utils"
)

type baseClient interface {
	DO(*http.Request) (*http.Response, error)
	BaseUrl() string
	Query() *query.Service
}

type Service struct {
	client baseClient
}

func NewService(client baseClient) *Service {
	return &Service{client: client}
}

func (s *Service) Create(input *CreateAlertRequest) (*Alert, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall create alert payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.client.BaseUrl()+"/api/v1/alert", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.DO(req)
	if err != nil {
		return nil, fmt.Errorf("alert creation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	var created Alert
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("failed to decode alert creation response: %w", err)
	}

	return &created, nil
}

func (s *Service) Get(id string) (*Alert, error) {
	url := fmt.Sprintf("%s/api/v1/alert/%s", s.client.BaseUrl(), id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build GET request: %w", err)
	}

	resp, err := s.client.DO(req)
	if err != nil {
		return nil, fmt.Errorf("alert get request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	var alert Alert
	if err := json.NewDecoder(resp.Body).Decode(&alert); err != nil {
		return nil, fmt.Errorf("failed to decode alert response: %w", err)
	}

	return &alert, nil
}

func (s *Service) Delete(id string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/alert/%s", s.client.BaseUrl(), id), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := s.client.DO(req)
	if err != nil {
		return fmt.Errorf("failed to execute delete request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		// Success
		return nil
	case http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusInternalServerError:
		return fmt.Errorf("delete failed with status %d", resp.StatusCode)
	default:
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}
}

func (s *Service) Update(id string, input *UpdateAlertRequest) error {
	payload, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshall update alert payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/alert/%s", s.client.BaseUrl(), id)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.DO(req)
	if err != nil {
		return fmt.Errorf("alert update request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return nil
}

func (s *Service) List(options ...query.Option) ([]Alert, error) {
	options = utils.Prepend(options, query.WithListing("listAlert"))

	responseBody, err := s.client.Query().Execute(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	var alerts []Alert
	if err := json.Unmarshal(responseBody, &alerts); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return alerts, nil
}
