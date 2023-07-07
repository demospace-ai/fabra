package temporal

import (
	"context"

	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/server/common/views"
)

type UpdateCursorInput struct {
	Sync           views.Sync
	CursorPosition string
}

func (a *Activities) UpdateCursor(ctx context.Context, input UpdateCursorInput) error {
	sync, err := syncs.LoadSyncByID(a.Db, input.Sync.OrganizationID, input.Sync.ID)
	if err != nil {
		return errors.Wrap(err, "(temporal.UpdateCursorActivity) loading sync")
	}

	_, err = syncs.UpdateCursor(a.Db, sync, input.CursorPosition)
	if err != nil {
		return errors.Wrap(err, "(temporal.UpdateCursorActivity) updating cursor")
	}

	return nil
}
