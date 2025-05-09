package scheduler

// Publishers are used on trigger conditions
// For now, only an HTTP publisher is present which is used for posting events to the webserver
// In the future, this should be swapped to a message queue between webserver and worker.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPEventPublisher publishes events to an HTTP endpoint
type HTTPEventPublisher struct {
	eventURL   string
	httpClient *http.Client
}

func NewHTTPEventPublisher(eventURL string) *HTTPEventPublisher {
	return &HTTPEventPublisher{
		eventURL:   eventURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *HTTPEventPublisher) PublishEvent(ctx context.Context, eventData map[string]any) error {
	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.eventURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	return nil
}
