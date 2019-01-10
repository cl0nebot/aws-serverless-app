package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func textRequest(value string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body:    value,
	}
}

func jsonRequest(org *OrgStatus) events.APIGatewayProxyRequest {
	body, _ := json.Marshal(org)
	return events.APIGatewayProxyRequest{
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:    string(body),
	}
}

func TestJsonRequestHandler(t *testing.T) {
	org := OrgStatus{
		OrgID:         "org1",
		Status:        "Active",
		EffectiveDate: "2018-12-12",
	}

	t.Run("Inforce JSON Request", func(t *testing.T) {
		resp, err := handler(jsonRequest(&org))
		if err != nil {
			t.Fatalf("Failed with error %+v", err)
		}
		var respOrg OrgStatus
		err = json.Unmarshal([]byte(resp.Body), &respOrg)
		if err != nil {
			t.Fatalf("Invalid JSON response %s", resp.Body)
		}
		if !respOrg.Inforce {
			t.Fatal("org should be in-force")
		}
	})

	t.Run("Inactive JSON Request", func(t *testing.T) {
		org.Status = "Inactive"
		resp, err := handler(jsonRequest(&org))
		if err != nil {
			t.Fatalf("Failed with error %+v", err)
		}
		var respOrg OrgStatus
		err = json.Unmarshal([]byte(resp.Body), &respOrg)
		if err != nil {
			t.Fatalf("Invalid JSON response %s", resp.Body)
		}
		if respOrg.Inforce {
			t.Fatal("org should not be in-force")
		}
	})

	t.Run("In-effective JSON Request", func(t *testing.T) {
		org.Status = "Active"
		org.EffectiveDate = "2020-01-01"
		resp, err := handler(jsonRequest(&org))
		if err != nil {
			t.Fatalf("Failed with error %+v", err)
		}
		var respOrg OrgStatus
		err = json.Unmarshal([]byte(resp.Body), &respOrg)
		if err != nil {
			t.Fatalf("Invalid JSON response %s", resp.Body)
		}
		if respOrg.Inforce {
			t.Fatal("org should not be in-force")
		}
	})

	t.Run("Bad JSON Request", func(t *testing.T) {
		resp, err := handler(events.APIGatewayProxyRequest{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    "BadRequest",
		})
		if err == nil {
			t.Fatal("Bad request should faile with error")
		}
		if resp.StatusCode != 400 {
			t.Fatalf("Bad request status code %d should be 400", resp.StatusCode)
		}
	})
}

func TestTextRequestHandler(t *testing.T) {

	t.Run("Inforce Text Request", func(t *testing.T) {
		resp, err := handler(textRequest("org1,Active,2018-12-01"))
		if err != nil {
			t.Fatalf("Failed with error %+v", err)
		}
		if resp.Body != "true" {
			t.Fatal("org should be in-force")
		}
	})

	t.Run("Inactive Text Request", func(t *testing.T) {
		resp, err := handler(textRequest("org1,Inactive,2018-12-01"))
		if err != nil {
			t.Fatalf("Failed with error %+v", err)
		}
		if resp.Body != "false" {
			t.Fatal("org should not be in-force")
		}
	})

	t.Run("In-effective Text Request", func(t *testing.T) {
		resp, err := handler(textRequest("org1,Active,2020-01-01"))
		if err != nil {
			t.Fatalf("Failed with error %+v", err)
		}
		if resp.Body != "false" {
			t.Fatal("org should not be in-force")
		}
	})

	t.Run("Bad Text Request", func(t *testing.T) {
		resp, err := handler(events.APIGatewayProxyRequest{
			Headers: map[string]string{"Content-Type": "text/plain"},
			Body:    "",
		})
		if err == nil {
			t.Fatal("Bad request should faile with error")
		}
		if resp.StatusCode != 400 {
			t.Fatalf("Bad request status code %d should be 400", resp.StatusCode)
		}
	})
}
