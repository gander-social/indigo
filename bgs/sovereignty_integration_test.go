package bgs

import (
	"context"
	"testing"
	"time"

	comatproto "github.com/gander-social/gander-indigo-sovereign/api/atproto"
	"github.com/gander-social/gander-indigo-sovereign/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper to create BGS instance for testing
func createTestBGS() *BGS {
	config := &BGSConfig{
		SovereigntyEnabled:   true,
		SovereignCountryCode: "CA",
		SovereigntyMode:      "strict",
	}

	return &BGS{
		config:             config,
		geographicFilter:   NewCanadianGeographicFilter(),
		sovereigntyMetrics: NewSovereigntyMetrics(), // Uses singleton
	}
}

func TestBGSSovereignEventFilter_Integration(t *testing.T) {
	// Create a BGS with sovereignty enabled
	bgs := createTestBGS()

	// Initialize the filter
	ctx := context.Background()
	require.NoError(t, bgs.geographicFilter.Initialize(ctx))

	// Test Canadian event passes through
	canadianEvent := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:toronto-user-123",
			Seq:  1,
			Time: time.Now().Format(time.RFC3339),
		},
	}

	result := bgs.sovereignEventFilter(canadianEvent)
	assert.True(t, result, "Canadian event should pass through sovereign filter")

	// Test non-Canadian event is filtered out
	nonCanadianEvent := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:american-user-123",
			Seq:  2,
			Time: time.Now().Format(time.RFC3339),
		},
	}

	result = bgs.sovereignEventFilter(nonCanadianEvent)
	assert.False(t, result, "Non-Canadian event should be filtered out")
}

func TestBGSSovereignEventFilter_NilFilter(t *testing.T) {
	// Create a BGS without geographic filter
	bgs := &BGS{
		config: &BGSConfig{
			SovereigntyEnabled: false,
		},
		geographicFilter: nil,
	}

	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:any-user",
			Seq:  1,
			Time: time.Now().Format(time.RFC3339),
		},
	}

	result := bgs.sovereignEventFilter(event)
	assert.False(t, result, "Events should be filtered out when no geographic filter is available")
}

func TestBGSSovereignEventFilter_IdentityEvents(t *testing.T) {
	// Create a BGS with sovereignty enabled
	bgs := createTestBGS()

	ctx := context.Background()
	require.NoError(t, bgs.geographicFilter.Initialize(ctx))

	// Test Canadian identity event
	identityEvent := &events.XRPCStreamEvent{
		RepoIdentity: &comatproto.SyncSubscribeRepos_Identity{
			Did: "did:plc:montreal-user-456",
			Seq: 1,
		},
	}

	result := bgs.sovereignEventFilter(identityEvent)
	assert.True(t, result, "Canadian identity event should pass through")

	// Test non-Canadian identity event
	nonCanadianIdentityEvent := &events.XRPCStreamEvent{
		RepoIdentity: &comatproto.SyncSubscribeRepos_Identity{
			Did: "did:plc:london-user-789",
			Seq: 2,
		},
	}

	result = bgs.sovereignEventFilter(nonCanadianIdentityEvent)
	assert.False(t, result, "Non-Canadian identity event should be filtered out")
}

func TestBGSSovereignEventFilter_AccountEvents(t *testing.T) {
	// Create a BGS with sovereignty enabled
	bgs := createTestBGS()

	ctx := context.Background()
	require.NoError(t, bgs.geographicFilter.Initialize(ctx))

	// Test Canadian account event
	active := true
	accountEvent := &events.XRPCStreamEvent{
		RepoAccount: &comatproto.SyncSubscribeRepos_Account{
			Did:    "did:plc:vancouver-user-789",
			Seq:    1,
			Active: active,
		},
	}

	result := bgs.sovereignEventFilter(accountEvent)
	assert.True(t, result, "Canadian account event should pass through")
}

// Benchmark the sovereign event filter performance
func BenchmarkBGSSovereignEventFilter(b *testing.B) {
	bgs := createTestBGS()

	ctx := context.Background()
	bgs.geographicFilter.Initialize(ctx)

	event := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:toronto-benchmark-user",
			Seq:  1,
			Time: time.Now().Format(time.RFC3339),
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bgs.sovereignEventFilter(event)
	}
}

func BenchmarkBGSSovereignEventFilter_Mixed(b *testing.B) {
	bgs := createTestBGS()

	ctx := context.Background()
	bgs.geographicFilter.Initialize(ctx)

	canadianEvent := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:toronto-user",
			Seq:  1,
			Time: time.Now().Format(time.RFC3339),
		},
	}

	nonCanadianEvent := &events.XRPCStreamEvent{
		RepoCommit: &comatproto.SyncSubscribeRepos_Commit{
			Repo: "did:plc:american-user",
			Seq:  2,
			Time: time.Now().Format(time.RFC3339),
		},
	}

	events := []*events.XRPCStreamEvent{canadianEvent, nonCanadianEvent}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, evt := range events {
			bgs.sovereignEventFilter(evt)
		}
	}
}
