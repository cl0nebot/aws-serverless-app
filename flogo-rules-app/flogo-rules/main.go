package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// OrgStatus from request
	var reqBody OrgStatus

	log.Printf("Received request %+v\n", request)
	reqType := request.Headers["Content-Type"]
	log.Printf("request Content-Type: %s\n", reqType)
	if reqType == "" {
		// API gateway renames header name Content-Type to lower case
		reqType = request.Headers["content-type"]
		log.Printf("request content-type: %s\n", reqType)
	}
	if reqType == "text/plain" {
		// assume csv data 'OrgID,Status,EffectiveDate'
		tokens := strings.Split(request.Body, ",")
		if len(tokens) != 3 {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
				Body: "{}",
			}, fmt.Errorf("Invalid request %s. it should be of format 'orgID,statu,effectiveDate'", request.Body)
		}
		reqBody = OrgStatus{
			OrgID:         tokens[0],
			Status:        tokens[1],
			EffectiveDate: tokens[2],
		}
	} else {
		// default content-type is application/json, so parse json request body
		reqType = "application/json"
		if err := json.Unmarshal([]byte(request.Body), &reqBody); err != nil {
			// return client request error
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: "{}",
			}, err
		}
	}

	if reqBody.OrgID == "" {
		// return client request error
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": reqType,
			},
			Body: "{}",
		}, errors.New("Bad request: missing org-ID")
	}

	// invoke orgstatus rules
	invokeRules(&reqBody)

	// return the result
	result := strconv.FormatBool(reqBody.Inforce)
	if reqType == "application/json" {
		result = fmt.Sprintf("{\"inforce\": %t}", reqBody.Inforce)
	}
	log.Printf("Return Inforce: %s", result)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": reqType,
		},
		Body: result,
	}, nil
}

func main() {
	lambda.Start(handler)
}
