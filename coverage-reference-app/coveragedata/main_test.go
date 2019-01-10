package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("Invalid PUT Request", func(t *testing.T) {
		respEvt, err := handler(events.APIGatewayProxyRequest{
			Body: "BadRequest",
		})
		if err == nil {
			t.Fatal("Bad PUT request should generate error")
		}
		if respEvt.StatusCode != 400 {
			t.Fatal("Bad PUT request should return status 400")
		}
	})
}
