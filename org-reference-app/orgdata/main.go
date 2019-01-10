package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("APIGatewayProxyRequest %+v\n", request)
	var result string
	var err error

	startTime := time.Now()
	if request.HTTPMethod == "GET" {
		// query ?key=org1, returns value of the specified key from Redis
		key := request.QueryStringParameters["key"]
		log.Printf("query parameter for key %s\n", key)
		result, err = getFromRedis(key)
	} else {
		// assume POST with format key=value or init=#
		tokens := strings.Split(request.Body, "=")
		if len(tokens) != 2 {
			log.Printf("Error: invalid PUT request %s, it should be of format key=value", request.Body)
			return events.APIGatewayProxyResponse{
				Body:       "Bad Request",
				StatusCode: 400,
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, fmt.Errorf("Invalid put request: %s", request.Body)
		}
		if tokens[0] == "init" {
			// initialize Redis by inserting specified number of random keys
			result = fmt.Sprintf("Initialized %s keys", tokens[1])
			var n int
			if n, err = strconv.Atoi(tokens[1]); err == nil {
				err = initKeys(n)
			}
		} else {
			result = fmt.Sprintf("Put key %s = %s", tokens[0], tokens[1])
			err = putToRedis(tokens[0], tokens[1])
		}
	}
	log.Printf("Redis request elapsed time %s\n", time.Since(startTime))

	if err != nil {
		log.Printf("Error: %s - %+v", result, err)
		return events.APIGatewayProxyResponse{
			Body:       result,
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       result,
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
