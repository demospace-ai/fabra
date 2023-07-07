package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/input"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/views"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.fabra.io/server/common/test"
	"go.fabra.io/server/internal/api"
)

var _ = Describe("Sending an ObjectField batch update request", func() {
	var auth auth.Authentication
	var object *models.Object
	var makeRequest func(body interface{}) *http.Request

	BeforeEach(func() {
		auth = getAuth(db)
		destination, _ := test.CreateDestination(db, auth.Organization.ID)
		object = test.CreateObject(db, auth.Organization.ID, destination.ID, models.SyncModeFullOverwrite)
		makeRequest = func(body interface{}) *http.Request {
			jsonBody, _ := json.Marshal(body)
			request := httptest.NewRequest("POST", fmt.Sprintf("/object/%d/object_fields", object.ID), bytes.NewReader(jsonBody))
			return mux.SetURLVars(request, map[string]string{
				"objectID": fmt.Sprintf("%d", object.ID),
			})
		}
	})

	Context("with an empty list", func() {
		It("should return a 200 status code", func() {
			response := httptest.NewRecorder()
			request := makeRequest(map[string]interface{}{
				"object_fields": []interface{}{},
			})
			err := service.UpdateObjectFields(auth, response, request)
			Expect(err).To(BeNil(), "no error should be returned, got %s", err)
			Expect(response.Code).To(Equal(200))
		})
	})

	Context("with an object body that's missing ID", func() {
		It("should fail validation", func() {
			response := httptest.NewRecorder()
			request := makeRequest(map[string]interface{}{
				"object_fields": []map[string]interface{}{
					{}, // Missing ID
				},
			})
			err := service.UpdateObjectFields(auth, response, request)
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			fieldError := err.(validator.ValidationErrors)[0]
			Expect(fieldError.Field()).To(Equal("ID"))
		})
	})

	Context("with an object id but no change", func() {
		It("should return a 200 status code", func() {
			desc := "test description (shouldn't change)"
			objFields := test.CreateObjectFields(db, object.ID, []input.ObjectField{
				{
					Name:        "test (shouldn't change)",
					Description: &desc, // Description will not be updated
				},
			})
			response := httptest.NewRecorder()
			request := makeRequest(map[string]interface{}{
				"object_fields": []map[string]interface{}{
					{
						"id": objFields[0].ID,
						// Do not provide Description (This tests partial update)
					},
				},
			})
			err := service.UpdateObjectFields(auth, response, request)
			Expect(err).To(BeNil(), "no error should be returned, got %s", err)
			Expect(response.Code).To(Equal(200))
			var actual api.CreateObjectFieldsResponse
			json.Unmarshal(response.Body.Bytes(), &actual)
			Expect(*actual.ObjectFields[0].Description).To(Equal("test description (shouldn't change)"))
			Expect(actual.ObjectFields[0].Name).To(Equal("test (shouldn't change)"))
		})
	})

	Context("to change an object's properties", func() {
		It("should return a 200 status code, and update the anme", func() {
			disname := "old display name"
			objField := test.CreateObjectFields(db, object.ID, []input.ObjectField{
				{
					Name:        "old name",
					Type:        data.FieldTypeString,
					Description: nil,      // description will be updated from null to "new description"
					DisplayName: &disname, // display name will be updated from "old display name" to null
				},
			})[0]
			response := httptest.NewRecorder()
			desc := "new description"
			request := makeRequest(map[string]interface{}{
				"object_fields": []map[string]interface{}{
					{
						"id":           objField.ID,
						"name":         "new name", // This should be ignored because we don't allow updating the name
						"type":         "integer",  // This should be ignored because we don't allow updating the type
						"description":  desc,
						"display_name": nil, // This will set {"display_name": null}
					},
				},
			})
			err := service.UpdateObjectFields(auth, response, request)
			Expect(err).To(BeNil(), "no error should be returned, got %s", err)
			Expect(response.Code).To(Equal(200))
			expected, _ := json.Marshal(api.UpdateObjectFieldsResponse{
				ObjectFields: []views.ObjectField{
					{
						ID:          objField.ID,
						Name:        "old name",
						Type:        data.FieldTypeString,
						Description: &desc,
						DisplayName: nil, // Expects {"display_name": null} (or no display_name key)
					},
				},
				Failures: []int64{},
			})
			Expect(response.Body).To(MatchJSON(expected))
		})
	})

	Context("with an object id that doesn't exist", func() {
		It("should return a 200 status code, include the id in failures", func() {
			response := httptest.NewRecorder()
			request := makeRequest(map[string]interface{}{
				"object_fields": []map[string]interface{}{
					{
						"id": math.MaxInt64,
					},
				},
			})
			err := service.UpdateObjectFields(auth, response, request)
			Expect(err).To(BeNil(), "no error should be returned, got %s", err)
			Expect(response.Code).To(Equal(200))
			expect, _ := json.Marshal(api.UpdateObjectFieldsResponse{
				ObjectFields: []views.ObjectField{},
				Failures: []int64{
					math.MaxInt64,
				},
			})
			Expect(response.Body).To(MatchJSON(expect))
		})
	})
})
