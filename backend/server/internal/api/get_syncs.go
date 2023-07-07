package api

import (
	"encoding/json"
	"net/http"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/objects"
	"go.fabra.io/server/common/repositories/sources"
	sync_repository "go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/server/common/views"
)

// TODO: use graphql for this
type GetSyncsResponse struct {
	Syncs   []views.Sync   `json:"syncs"`
	Sources []views.Source `json:"sources"`
	Objects []views.Object `json:"objects"`
}

func (s ApiService) GetSyncs(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	syncs, err := sync_repository.LoadAllSyncs(s.db, auth.Organization.ID)
	if err != nil {
		return errors.Wrap(err, "(api.GetSyncs)")
	}

	sourceIDset := make(map[int64]bool)
	objectIDset := make(map[int64]bool)
	for _, sync := range syncs {
		sourceIDset[sync.SourceID] = true
		objectIDset[sync.ObjectID] = true
	}

	var sourceIDs []int64
	for key := range sourceIDset {
		sourceIDs = append(sourceIDs, key)
	}

	var objectIDs []int64
	for key := range objectIDset {
		objectIDs = append(objectIDs, key)
	}

	sources, err := sources.LoadSourcesByIDs(s.db, auth.Organization.ID, sourceIDs)
	if err != nil {
		return errors.Wrap(err, "(api.GetSyncs)")
	}

	objects, err := objects.LoadObjectsByIDs(s.db, auth.Organization.ID, objectIDs)
	if err != nil {
		return errors.Wrap(err, "(api.GetSyncs)")
	}

	syncViews := []views.Sync{}
	for _, sync := range syncs {
		syncViews = append(syncViews, views.ConvertSync(&sync))
	}

	objectViews := []views.Object{}
	for _, object := range objects {
		// Don't bother to fetch the object fields
		objectViews = append(objectViews, views.ConvertObject(&object, []models.ObjectField{}))
	}

	return json.NewEncoder(w).Encode(GetSyncsResponse{
		Syncs:   syncViews,
		Sources: views.ConvertSourceConnections(sources),
		Objects: objectViews,
	})
}
