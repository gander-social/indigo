package bgs

import (
	"time"
)

// StreamEvent represents a unified event type for geographic filtering
type StreamEvent struct {
	Repo     string    `json:"repo"`
	Kind     string    `json:"kind"`
	Time     time.Time `json:"time"`
	Seq      int64     `json:"seq"`
	Source   string    `json:"source,omitempty"`   // Source PDS endpoint
	Country  string    `json:"country,omitempty"`  // Detected country
	Verified bool      `json:"verified,omitempty"` // Whether location is verified
}

// NewStreamEvent creates a new StreamEvent
func NewStreamEvent(repo, kind string) *StreamEvent {
	return &StreamEvent{
		Repo: repo,
		Kind: kind,
		Time: time.Now(),
	}
}
