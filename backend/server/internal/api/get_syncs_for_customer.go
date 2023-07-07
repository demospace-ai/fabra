package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/repositories/objects"
	"go.fabra.io/server/common/repositories/sources"
	sync_repository "go.fabra.io/server/common/repositories/syncs"
	"go.fabra.io/server/common/views"
)

func (s ApiService) GetSyncsForCustomer(auth auth.Authentication, w http.ResponseWriter, r *http.Request) error {
	if auth.Organization == nil {
		return errors.Wrap(errors.NewBadRequest("must setup organization first"), "(api.GetSyncsForCustomer)")
	}

	vars := mux.Vars(r)
	endCustomerId, ok := vars["endCustomerId"]
	if !ok {
		return errors.Newf("(api.GetSyncsForCustomer) missing end customer ID from GetSyncsForCustomer request URL: %s", r.URL.RequestURI())
	}

	syncs, sources, objects, err := s.getSyncsForCustomer(auth, endCustomerId)
	if err != nil {
		return errors.Wrap(err, "(api.GetSyncsForCustomer)")
	}

	return json.NewEncoder(w).Encode(GetSyncsResponse{
		Syncs:   syncs,
		Sources: sources,
		Objects: objects,
	})
}

func (s ApiService) getSyncsForCustomer(auth auth.Authentication, endCustomerID string) ([]views.Sync, []views.Source, []views.Object, error) {

	syncs, err := sync_repository.LoadAllSyncsForCustomer(s.db, auth.Organization.ID, auth.LinkToken.EndCustomerID)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "(api.getSyncsForCustomer)")
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

	sources, err := sources.LoadSourcesByIDsForCustomer(s.db, auth.Organization.ID, auth.LinkToken.EndCustomerID, sourceIDs)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "(api.getSyncsForCustomer)")
	}

	objects, err := objects.LoadObjectsByIDs(s.db, auth.Organization.ID, objectIDs)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "(api.getSyncsForCustomer)")
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

	return syncViews, views.ConvertSourceConnections(sources), objectViews, nil
}
