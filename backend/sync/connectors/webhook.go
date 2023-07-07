package connectors

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"go.fabra.io/server/common/crypto"
	"go.fabra.io/server/common/data"
	"go.fabra.io/server/common/errors"
	"go.fabra.io/server/common/query"
	"go.fabra.io/server/common/views"
	"golang.org/x/time/rate"
)

const MAX_WEBHOOK_BATCH_SIZE = 1_000
const REFILL_RATE = 100
const MAX_BURST = 100

type WebhookData struct {
	ObjectID          int64            `json:"object_id"`
	ObjectName        string           `json:"object_name"`
	EndCustomerID     string           `json:"end_customer_id"`
	EndCustomerApiKey *string          `json:"end_customer_api_key,omitempty"`
	FabraTimestamp    int64            `json:"fabra_timestamp"`
	Data              []map[string]any `json:"data"`
}

type WebhookImpl struct {
	queryService               query.QueryService
	cryptoService              crypto.CryptoService
	encryptedEndCustomerApiKey *string
}

func NewWebhookConnector(queryService query.QueryService, cryptoService crypto.CryptoService, encryptedEndCustomerApiKey *string) Connector {
	return WebhookImpl{
		queryService:               queryService,
		cryptoService:              cryptoService,
		encryptedEndCustomerApiKey: encryptedEndCustomerApiKey, // TODO: does this belong here?
	}
}

func (wh WebhookImpl) Read(
	ctx context.Context,
	sourceConnection views.FullConnection,
	sync views.Sync,
	fieldMappings []views.FieldMapping,
	rowsC chan<- []data.Row,
	readOutputC chan<- ReadOutput,
	errC chan<- error,
) {
	errC <- errors.New("webhook source not supported")
}

func (wh WebhookImpl) Write(
	ctx context.Context,
	destinationConnection views.FullConnection,
	destinationOptions DestinationOptions,
	object views.Object,
	sync views.Sync,
	fieldMappings []views.FieldMapping,
	rowsC <-chan []data.Row,
	writeOutputC chan<- WriteOutput,
	errC chan<- error,
) {
	// TODO: allow customizing the rate limit
	limiter := rate.NewLimiter(REFILL_RATE, MAX_BURST)

	decryptedSigningKey, err := wh.cryptoService.DecryptWebhookSigningKey(destinationConnection.Credentials)
	if err != nil {
		errC <- err
		return
	}

	var decryptedEndCustomerApiKey *string
	if wh.encryptedEndCustomerApiKey != nil {
		decryptedEndCustomerApiKey, err = wh.cryptoService.DecryptEndCustomerApiKey(*wh.encryptedEndCustomerApiKey)
		if err != nil {
			errC <- err
			return
		}
	}

	orderedObjectFields := wh.createOrderedObjectFields(object.ObjectFields, fieldMappings)
	outputDataList := []map[string]any{}

	rowsWritten := 0
	for {
		currentBatchSize := 0
		rows, more := <-rowsC
		if !more {
			break
		}

		rowsWritten += len(rows)
		for _, row := range rows {
			outputData := map[string]any{}
			for i, value := range row {
				fieldMapping := fieldMappings[i]
				destFieldName := orderedObjectFields[i].Name
				// add raw values to the json object even if they're nil
				if fieldMapping.IsJsonField {
					existing, ok := outputData[destFieldName]
					if !ok {
						existing = make(map[string]any)
						outputData[destFieldName] = existing
					}

					existing.(map[string]any)[fieldMapping.SourceFieldName] = value
				} else {
					if value != nil {
						outputData[destFieldName] = value
					}
				}
			}
			outputDataList = append(outputDataList, outputData)

			currentBatchSize++
			// TODO: allow customizing batch size
			if currentBatchSize == MAX_WEBHOOK_BATCH_SIZE {
				// TODO: add retry
				limiter.Wait(ctx)
				err := wh.sendData(object, sync.EndCustomerID, decryptedEndCustomerApiKey, outputDataList, destinationConnection.Host, *decryptedSigningKey)
				if err != nil {
					errC <- err
					return
				}

				currentBatchSize = 0
				outputDataList = nil
			}
		}

		if currentBatchSize > 0 {
			err := wh.sendData(object, sync.EndCustomerID, decryptedEndCustomerApiKey, outputDataList, destinationConnection.Host, *decryptedSigningKey)
			if err != nil {
				errC <- err
				return
			}
		}
	}

	writeOutputC <- WriteOutput{
		RowsWritten: rowsWritten,
	}

	close(errC)
}

func (wh WebhookImpl) sendData(object views.Object, endCustomerID string, endCustomerApiKey *string, outputDataList []map[string]any, webhookUrl string, decryptedSigningKey string) error {
	webhookData := WebhookData{
		ObjectID:          object.ID,
		ObjectName:        object.DisplayName,
		EndCustomerID:     endCustomerID,
		EndCustomerApiKey: endCustomerApiKey,
		FabraTimestamp:    time.Now().Unix(),
		Data:              outputDataList,
	}
	marshalled, err := json.Marshal(webhookData)
	if err != nil {
		return errors.Wrap(err, "(connectors.WebhookImpl.sendData)")
	}

	request, _ := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(marshalled))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("X-FABRA-SIGNATURE", wh.signPayload(decryptedSigningKey, marshalled))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "(connectors.WebhookImpl.sendData)")
	}
	response.Body.Close()

	return nil
}

func (wh WebhookImpl) signPayload(secret string, data []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func (wh WebhookImpl) createOrderedObjectFields(objectFields []views.ObjectField, fieldMappings []views.FieldMapping) []views.ObjectField {
	objectFieldIdToObjectField := make(map[int64]views.ObjectField)
	for _, objectField := range objectFields {
		objectFieldIdToObjectField[objectField.ID] = objectField
	}

	var orderedObjectFields []views.ObjectField
	for _, fieldMapping := range fieldMappings {
		orderedObjectFields = append(orderedObjectFields, objectFieldIdToObjectField[fieldMapping.DestinationFieldId])
	}

	return orderedObjectFields
}
