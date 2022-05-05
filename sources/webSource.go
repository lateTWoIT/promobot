package sources

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type WebSource struct {
	name string
	url  string
}

func NewWebSource(name, url string) *WebSource {
	return &WebSource{
		name: name,
		url:  url,
	}
}

func (c *WebSource) Name() string { return c.name }
func (c *WebSource) Type() string { return "Web" }

func (c *WebSource) Retrieve(loc *time.Location) ([]Event, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return []Event{}, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Event{}, err
	}

	var WebResp WebResponse
	if err = json.Unmarshal(respBytes, &WebResp); err != nil {
		return []Event{}, err
	}

	events := make([]Event, 0)
	for _, evt := range WebResp.Results {
		if haveNextInSeries(events, evt.Name) {
			continue
		}
		events = append(events, Event{
			Name:     evt.Name,
			Source:   c.name,
			URL:      evt.URL,
			Location: evt.Venue.Name,
			DateTime: time.UnixMilli(evt.Time).In(loc),
		})
	}

	return events, nil
}

func haveNextInSeries(events []Event, eventName string) bool {
	for _, e := range events {
		if e.Name == eventName {
			return true
		}
	}
	return false
}

type WebResponse struct {
	Results []WebEvent `json:"results"`
}

type WebEvent struct {
	Name  string      `json:"name"`
	Time  int64       `json:"time"`
	URL   string      `json:"event_url"`
	Venue WebVenue `json:"venue"`
}

type WebVenue struct {
	Name string `json:"name"`
}
