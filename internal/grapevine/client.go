package grapevine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Urgency string

const (
	UrgencyVeryLow Urgency = "very-low"
	UrgencyLow     Urgency = "low"
	UrgencyNormal  Urgency = "normal"
	UrgencyHigh    Urgency = "high"
)

type Notification struct {
	TTL     int     `json:"ttl"`
	Urgency Urgency `json:"urgency"`
	Title   string  `json:"title"`
	Body    string  `json:"body"`
}

func SendNotification(ctx context.Context, endpoint string, topic string, notification Notification) error {
	body, err := json.Marshal(&notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/api/v1/notifications/"+url.PathEscape(topic), bytes.NewReader(body))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("got unexpected status code: %d", res.StatusCode)
	}

	return nil
}
