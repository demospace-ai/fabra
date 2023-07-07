package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"go.fabra.io/server/common/auth"
	"go.fabra.io/server/common/models"
	"go.fabra.io/server/common/test"
	"go.fabra.io/server/internal/api"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.fabra.io/server/common/data"
)

var _ = Describe("Sending an ObjectField creation request", func() {
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
			err := service.CreateObjectFields(auth, response, makeRequest(map[string]interface{}{
				"object_fields": []interface{}{},
			}))
			Expect(err).To(BeNil(), "no error should be returned, got %s", err)
			Expect(response.Code).To(Equal(200))
		})
	})

	Context("with an object that is {}", func() {
		It("should fail validation", func() {
			response := httptest.NewRecorder()
			err := service.CreateObjectFields(auth, response, makeRequest(map[string]interface{}{
				"object_fields": []interface{}{
					map[string]interface{}{
						// empty object
					},
				},
			}))
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			fieldError := err.(validator.ValidationErrors)[0]
			Expect(fieldError.Field()).To(Equal("Name"))
		})
	})

	Context("with an object", func() {
		It("should succeed", func() {
			response := httptest.NewRecorder()
			err := service.CreateObjectFields(auth, response, makeRequest(map[string]interface{}{
				"object_fields": []interface{}{
					map[string]interface{}{
						"name":       "test",
						"field_type": "STRING",
					},
				},
			}))
			Expect(err).To(BeNil(), "no error should be returned, got %s", err)
			Expect(response.Code).To(Equal(200))
			var actual api.CreateObjectFieldsResponse
			err = json.Unmarshal(response.Body.Bytes(), &actual)
			Expect(err).To(BeNil(), "no error should be returned, got %s", err)
			Expect(actual.ObjectFields).To(HaveLen(1))
			Expect(actual.ObjectFields[0].Name).To(Equal("test"))
			Expect(actual.ObjectFields[0].Type).To(Equal(data.FieldTypeString))
		})
	})
})
