package main

import (
	"github.com/gander-social/gander-indigo-sovereign/bgs"
	"github.com/gander-social/gander-indigo-sovereign/cmd/relay/stream"
	"time"
)

// Use StreamEvent from bgs package
type StreamEvent = bgs.StreamEvent

// ConvertXRPCToStreamEvent converts stream.XRPCStreamEvent to StreamEvent
func ConvertXRPCToStreamEvent(xrpc *stream.XRPCStreamEvent) *bgs.StreamEvent {
	if xrpc == nil {
		return nil
	}

	event := &bgs.StreamEvent{
		Time: time.Now(),
	}

	// Handle different event types
	if xrpc.RepoCommit != nil {
		event.Repo = xrpc.RepoCommit.Repo
		event.Kind = "commit"
		event.Seq = xrpc.RepoCommit.Seq
	} else if xrpc.RepoIdentity != nil {
		event.Repo = xrpc.RepoIdentity.Did
		event.Kind = "identity"
		event.Seq = xrpc.RepoIdentity.Seq
	} else if xrpc.RepoAccount != nil {
		event.Repo = xrpc.RepoAccount.Did
		event.Kind = "account"
		event.Seq = xrpc.RepoAccount.Seq
	}

	return event
}
