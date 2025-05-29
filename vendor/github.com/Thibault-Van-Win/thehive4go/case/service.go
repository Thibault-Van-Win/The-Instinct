package cases

import (
	"encoding/json"
	"fmt"
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

func (s *Service) List(options ...query.Option) ([]Case, error) {
	options = utils.Prepend(options, query.WithListing("listCase"))

	responseBody, err := s.client.Query().Execute(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	// Parse response
	var cases []Case
	if err := json.Unmarshal(responseBody, &cases); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return cases, nil
}
