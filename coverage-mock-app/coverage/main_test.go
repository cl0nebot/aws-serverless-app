package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func requestEvent(reqID, orgID, covID string) events.APIGatewayProxyRequest {
	req := EligibilityRequest{
		ResourceType: "EligibilityRequest",
		ID:           reqID,
		Patient:      ReferenceData{Reference: "deceased"},
		Organization: ReferenceData{Reference: orgID},
		Insurer:      ReferenceData{Reference: "cygna"},
		Coverage:     ReferenceData{Reference: covID},
	}
	body, _ := json.Marshal(&req)
	return events.APIGatewayProxyRequest{
		Body: string(body),
	}
}

func TestHandler(t *testing.T) {
	t.Run("Bad Request", func(t *testing.T) {
		respEvt, err := handler(events.APIGatewayProxyRequest{
			Body: "BadRequest",
		})
		if err == nil {
			t.Fatal("Error failed to trigger with a bad request")
		}
		if respEvt.StatusCode != 400 {
			t.Fatalf("Response statusCode %d should be 400", respEvt.StatusCode)
		}
	})

	t.Run("Successful Request", func(t *testing.T) {
		respEvt, err := handler(requestEvent("test-1", "provider-1", "coverage-1"))
		if err != nil {
			t.Fatal("Everything should be ok")
		}
		if respEvt.StatusCode != 200 {
			t.Fatalf("Response statusCode %d should be 200", respEvt.StatusCode)
		}
	})
}
