package bgs

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gander-social/gander-indigo-sovereign/events"
)

// GeographicFilter interface for content filtering
type GeographicFilter interface {
	ShouldIncludeInSovereignFeed(event *events.XRPCStreamEvent) bool
	Initialize(ctx context.Context) error
}

// CanadianGeographicFilter filters content for Canadian sovereignty
type CanadianGeographicFilter struct {
	cache      map[string]bool
	cacheMutex sync.RWMutex
	ttl        time.Duration
	knownDIDs  map[string]bool
}

// NewCanadianGeographicFilter creates a new Canadian geographic filter
func NewCanadianGeographicFilter() *CanadianGeographicFilter {
	return &CanadianGeographicFilter{
		cache:     make(map[string]bool),
		ttl:       24 * time.Hour,
		knownDIDs: make(map[string]bool),
	}
}

// NewGeographicFilter creates a geographic filter for any country
func NewGeographicFilter(countryCode string) *CanadianGeographicFilter {
	// For now, all filters use the same logic but could be extended
	filter := NewCanadianGeographicFilter()
	// Could customize behavior based on countryCode
	return filter
}

// ShouldIncludeInSovereignFeed determines if content should be in sovereign feed
func (f *CanadianGeographicFilter) ShouldIncludeInSovereignFeed(event *events.XRPCStreamEvent) bool {
	if event == nil {
		return false
	}

	// Extract the repo/DID from the event
	var repo string
	switch {
	case event.RepoCommit != nil:
		repo = event.RepoCommit.Repo
	case event.RepoIdentity != nil:
		repo = event.RepoIdentity.Did
	case event.RepoAccount != nil:
		repo = event.RepoAccount.Did
	case event.RepoSync != nil:
		repo = event.RepoSync.Repo
	default:
		return false // Unknown event type
	}

	return f.isLikelyCanadian(repo)
}

// Initialize sets up the filter
func (f *CanadianGeographicFilter) Initialize(ctx context.Context) error {
	// Initialize filter with any required setup
	return nil
}

// AddCanadianDID adds a known Canadian DID to the filter
func (f *CanadianGeographicFilter) AddCanadianDID(did string) {
	f.cacheMutex.Lock()
	defer f.cacheMutex.Unlock()
	f.knownDIDs[did] = true
	f.cache[did] = true
}

// GetCache returns the cache for testing purposes
func (f *CanadianGeographicFilter) GetCache() map[string]bool {
	f.cacheMutex.RLock()
	defer f.cacheMutex.RUnlock()

	result := make(map[string]bool)
	for k, v := range f.cache {
		result[k] = v
	}
	return result
}

// isLikelyCanadian uses heuristics to determine if a DID is Canadian
func (f *CanadianGeographicFilter) isLikelyCanadian(did string) bool {
	f.cacheMutex.RLock()
	if cached, exists := f.cache[did]; exists {
		f.cacheMutex.RUnlock()
		return cached
	}
	f.cacheMutex.RUnlock()

	// Check known DIDs first
	f.cacheMutex.RLock()
	if known, exists := f.knownDIDs[did]; exists {
		f.cacheMutex.RUnlock()
		f.setCached(did, known)
		return known
	}
	f.cacheMutex.RUnlock()

	// Use keyword-based heuristics
	canadianKeywords := []string{
		"canadian", "canada", "toronto", "vancouver", "montreal",
		"calgary", "ottawa", "edmonton", "winnipeg", "quebec",
		"halifax", "victoria", "saskatoon", "regina", "fredericton",
	}

	didLower := strings.ToLower(did)
	for _, keyword := range canadianKeywords {
		if strings.Contains(didLower, keyword) {
			f.setCached(did, true)
			return true
		}
	}

	f.setCached(did, false)
	return false
}

// setCached sets a value in the cache thread-safely
func (f *CanadianGeographicFilter) setCached(did string, value bool) {
	f.cacheMutex.Lock()
	defer f.cacheMutex.Unlock()
	f.cache[did] = value
}
