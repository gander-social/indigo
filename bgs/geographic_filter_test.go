package bgs

import (
	"context"
	"strings"
	"testing"

	comatproto "github.com/gander-social/gander-indigo-sovereign/api/atproto"
	"github.com/gander-social/gander-indigo-sovereign/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCanadianGeographicFilter(t *testing.T) {
	filter := NewCanadianGeographicFilter()
	require.NotNil(t, filter, "Filter should not be nil")
}

func TestGeographicFilter_Initialize(t *testing.T) {
	filter := NewCanadianGeographicFilter()
	ctx := context.Background()

	err := filter.Initialize(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, filter.GetCache())
}

func TestGeographicFilter_CanadianContent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Test Canadian DID
	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:canadian-test-user",
		},
	}

	result := filter.ShouldIncludeInSovereignFeed(event)
	assert.True(t, result, "Canadian content should be included")
}

func TestGeographicFilter_TorontoContent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:toronto-user-123",
		},
	}

	result := filter.ShouldIncludeInSovereignFeed(event)
	assert.True(t, result, "Toronto content should be included")
}

func TestGeographicFilter_NonCanadianContent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Test non-Canadian DID
	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:american-user",
		},
	}

	result := filter.ShouldIncludeInSovereignFeed(event)
	assert.False(t, result, "Non-Canadian content should be excluded")
}

func TestGeographicFilter_UnknownContent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Test unknown DID
	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:unknown789",
		},
	}

	result := filter.ShouldIncludeInSovereignFeed(event)
	assert.False(t, result, "Unknown DIDs should be excluded by default")
}

func TestGeographicFilter_IdentityEvent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Test Canadian identity event
	event := &events.XRPCStreamEvent{
		RepoIdentity: &comatproto.SyncSubscribeRepos_Identity{
			Did: "did:plc:montreal-identity-user",
		},
	}

	result := filter.ShouldIncludeInSovereignFeed(event)
	assert.True(t, result, "Canadian identity events should be included")
}

func TestGeographicFilter_AccountEvent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Test Canadian account event
	event := &events.XRPCStreamEvent{
		RepoAccount: &comatproto.SyncSubscribeRepos_Account{
			Did: "did:plc:vancouver-account-user",
		},
	}

	result := filter.ShouldIncludeInSovereignFeed(event)
	assert.True(t, result, "Canadian account events should be included")
}

func TestGeographicFilter_CountryConfiguration(t *testing.T) {
	// Test that filter can be configured for different countries
	caFilter := NewGeographicFilter("CA")
	usFilter := NewGeographicFilter("US")

	require.NotNil(t, caFilter, "CA filter should not be nil")
	require.NotNil(t, usFilter, "US filter should not be nil")

	// Test with neutral DID
	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:test123",
		},
	}

	// Both should exist and process events
	caResult := caFilter.ShouldIncludeInSovereignFeed(event)
	usResult := usFilter.ShouldIncludeInSovereignFeed(event)

	// Results may be the same for unknown DIDs
	assert.False(t, caResult, "Unknown DID should be excluded by CA filter")
	assert.False(t, usResult, "Unknown DID should be excluded by US filter")
}

func TestGeographicFilter_CacheOperations(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Test caching works
	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:toronto-cached",
		},
	}

	// First call should compute and cache
	result1 := filter.ShouldIncludeInSovereignFeed(event)
	assert.True(t, result1)

	// Second call should use cache
	result2 := filter.ShouldIncludeInSovereignFeed(event)
	assert.True(t, result2)
	assert.Equal(t, result1, result2)

	// Verify cache contains the entry
	cache := filter.GetCache()
	cached, exists := cache[event.RepoCommit.Repo]
	assert.True(t, exists, "DID should be cached")
	assert.True(t, cached, "Cached value should be true")
}

func TestGeographicFilter_NilEvent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	result := filter.ShouldIncludeInSovereignFeed(nil)
	assert.False(t, result, "Nil events should be excluded")
}

func TestGeographicFilter_EmptyEvent(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Event with no commit, identity, or account data
	event := &events.XRPCStreamEvent{}
	result := filter.ShouldIncludeInSovereignFeed(event)
	assert.False(t, result, "Empty events should be excluded")
}

func TestGeographicFilter_CanadianKeywords(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	testCases := []struct {
		repo     string
		expected bool
		name     string
	}{
		{"did:plc:vancouver-user", true, "Vancouver"},
		{"did:plc:montreal-test", true, "Montreal"},
		{"did:plc:calgary-person", true, "Calgary"},
		{"did:plc:ottawa-gov", true, "Ottawa"},
		{"did:plc:quebec-user", true, "Quebec"},
		{"did:plc:halifax-person", true, "Halifax"},
		{"did:plc:edmonton-user", true, "Edmonton"},
		{"did:plc:winnipeg-test", true, "Winnipeg"},
		{"did:plc:newyork-user", false, "New York"},
		{"did:plc:london-user", false, "London"},
		{"did:plc:paris-user", false, "Paris"},
		{"did:plc:tokyo-user", false, "Tokyo"},
		{"did:plc:sydney-user", false, "Sydney"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event := &events.XRPCStreamEvent{
				RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
					Repo: tc.repo,
				},
			}
			result := filter.ShouldIncludeInSovereignFeed(event)
			assert.Equal(t, tc.expected, result, "Failed for %s", tc.name)
		})
	}
}

func TestGeographicFilter_DifferentEventTypes(t *testing.T) {
	filter := NewCanadianGeographicFilter()
	repo := "did:plc:toronto-test-all-events"

	// Test RepoCommit event
	commitEvent := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: repo,
		},
	}
	assert.True(t, filter.ShouldIncludeInSovereignFeed(commitEvent), "Commit event should be included")

	// Test RepoIdentity event
	identityEvent := &events.XRPCStreamEvent{
		RepoIdentity: &comatproto.SyncSubscribeRepos_Identity{
			Did: repo,
		},
	}
	assert.True(t, filter.ShouldIncludeInSovereignFeed(identityEvent), "Identity event should be included")

	// Test RepoAccount event
	accountEvent := &events.XRPCStreamEvent{
		RepoAccount: &comatproto.SyncSubscribeRepos_Account{
			Did: repo,
		},
	}
	assert.True(t, filter.ShouldIncludeInSovereignFeed(accountEvent), "Account event should be included")
}

func TestGeographicFilter_Performance(t *testing.T) {
	filter := NewCanadianGeographicFilter()

	// Pre-populate cache with some entries
	testDIDs := []string{
		"did:plc:toronto-perf-1",
		"did:plc:vancouver-perf-2",
		"did:plc:montreal-perf-3",
		"did:plc:american-perf-4",
		"did:plc:unknown-perf-5",
	}

	events := make([]*events.XRPCStreamEvent, len(testDIDs))
	for i, did := range testDIDs {
		events[i] = &events.XRPCStreamEvent{
			RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
				Repo: did,
			},
		}
	}

	// First pass - populate cache
	for _, event := range events {
		filter.ShouldIncludeInSovereignFeed(event)
	}

	// Second pass - should be faster due to caching
	for _, event := range events {
		result := filter.ShouldIncludeInSovereignFeed(event)
		// Verify expected results
		if strings.Contains(event.RepoCommit.Repo, "toronto") ||
			strings.Contains(event.RepoCommit.Repo, "vancouver") ||
			strings.Contains(event.RepoCommit.Repo, "montreal") {
			assert.True(t, result, "Canadian DID should be included: %s", event.RepoCommit.Repo)
		} else {
			assert.False(t, result, "Non-Canadian DID should be excluded: %s", event.RepoCommit.Repo)
		}
	}
}
