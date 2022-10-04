package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charlieegan3/activities-rss/pkg/download"
)

// NewEntry is a job that creates a new entry in the feed with the activities from yesterday
type NewEntry struct {
	ScheduleOverride string
	Endpoint         string

	Host      string
	AthleteID string
	Email     string
	Password  string
}

func (n *NewEntry) Name() string {
	return "new-entry"
}

func (n *NewEntry) Run(ctx context.Context) error {
	doneCh := make(chan bool)
	errCh := make(chan error)

	go func() {
		activities, err := download.Download(n.Host, n.Email, n.Password, n.AthleteID)
		if err != nil {
			errCh <- fmt.Errorf("failed to download activities: %w", err)
			return
		}

		var date, html string
		for _, act := range activities {
			if !act.IsOnDate(time.Now().Add(-24 * time.Hour)) {
				continue // only report yesterday
			}

			html += act.ToHTML(n.Host)
			date = act.DateString()
		}

		datab := []map[string]string{
			{
				"title": "Activities on " + date,
				"body":  html,
				"url":   "",
			},
		}

		b, err := json.Marshal(datab)
		if err != nil {
			errCh <- fmt.Errorf("failed to form new item JSON: %s", err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", n.Endpoint, bytes.NewBuffer(b))
		if err != nil {
			errCh <- fmt.Errorf("failed to build request for new item: %s", err)
			return
		}

		req.Header.Add("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			errCh <- fmt.Errorf("failed to send request for new item: %s", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("failed to send request: non 200OK response")
			return
		}

		doneCh <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case e := <-errCh:
		return fmt.Errorf("job failed with error: %s", e)
	case <-doneCh:
		return nil
	}
}

func (n *NewEntry) Timeout() time.Duration {
	return 30 * time.Second
}

func (n *NewEntry) Schedule() string {
	if n.ScheduleOverride != "" {
		return n.ScheduleOverride
	}
	return "0 0 6 * * *"
}
