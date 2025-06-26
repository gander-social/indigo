// Copied from indigo:api/gndr/actorputPreferences.go

package agnostic

// schema: app.gndr.actor.putPreferences

import (
	"context"

	"github.com/gander-social/gander-indigo-sovereign/lex/util"
)

// ActorPutPreferences_Input is the input argument to a app.gndr.actor.putPreferences call.
type ActorPutPreferences_Input struct {
	Preferences []map[string]any `json:"preferences" cborgen:"preferences"`
}

// ActorPutPreferences calls the XRPC method "app.gndr.actor.putPreferences".
func ActorPutPreferences(ctx context.Context, c util.LexClient, input *ActorPutPreferences_Input) error {
	if err := c.LexDo(ctx, util.Procedure, "application/json", "app.gndr.actor.putPreferences", nil, input, nil); err != nil {
		return err
	}

	return nil
}
