package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var client *lambda.Lambda
var region, lambdaFunc string
var threads, countPerThread, targetMs, printInterval int
var debug bool
var perf chan int

func init() {
	flag.StringVar(&region, "region", "us-west-2", "AWS region")
	flag.StringVar(&lambdaFunc, "lambda", "arn:aws:lambda:us-west-2:742759186184:function:orchestrator-app-OrchestratorFunction-A4X63IX1EEDQ", "ARN of lambda function to test")
	flag.IntVar(&threads, "threads", 2, "Number of concurrent test threads")
	flag.IntVar(&countPerThread, "count", 10, "Number of requests per thread")
	flag.IntVar(&targetMs, "target", 100, "Target max elapsed time in millis")
	flag.IntVar(&printInterval, "print", 5, "Print stats when this number of tests are complete")
	flag.BoolVar(&debug, "debug", false, "Show debug logs")

	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client = lambda.New(sess, &aws.Config{Region: &region})
}

func main() {
	// parse command-line args
	flag.Parse()

	// channel for collecting performance data
	perf = make(chan int, 100)
	go collectStats()

	comp := make(chan int)
	for t := 0; t < threads; t++ {
		go func(id int) {
			for i := 0; i < countPerThread; i++ {
				invokeFunc(i)
			}
			comp <- id
		}(t)
	}

	// wait for all threads
	for t := 0; t < threads; t++ {
		id := <-comp
		fmt.Printf("Thread %d completed\n", id)
	}

	// send msg to quit stat collection, and wait 2 seconds
	perf <- -1
	time.Sleep(2 * time.Second)
}

// ReferenceData data
type ReferenceData struct {
	Reference string `json:"reference"`
}

// EligibilityRequest data
type EligibilityRequest struct {
	ResourceType string        `json:"resourceType"`
	ID           string        `json:"ID"`
	Status       string        `json:"status,omitempty"`
	Patient      ReferenceData `json:"patient,omitempty"`
	Created      string        `json:"created,omitempty"`
	Organization ReferenceData `json:"organization,omitempty"`
	Insurer      ReferenceData `json:"insurer,omitempty"`
	Coverage     ReferenceData `json:"coverage,omitempty"`
}

// RequestEvent for orchestrator lambda function
type RequestEvent struct {
	Body string
}

func invokeFunc(index int) {
	// construct orchestrator request
	reqID := fmt.Sprintf("TEST-%06d", index)
	covID := fmt.Sprintf("C-%06d", index%50)
	orgID := fmt.Sprintf("P-%06d", (index/50)%50)
	req := EligibilityRequest{
		ResourceType: "EligibilityRequest",
		ID:           reqID,
		Patient:      ReferenceData{Reference: "deceased"},
		Organization: ReferenceData{Reference: orgID},
		Insurer:      ReferenceData{Reference: "cygna"},
		Coverage:     ReferenceData{Reference: covID},
	}
	body, _ := json.Marshal(&req)
	evt := RequestEvent{
		Body: string(body),
	}

	payload, _ := json.Marshal(&evt)
	if debug {
		log.Printf("Send request to orchestrator lambda %s: %s\n", lambdaFunc, string(payload))
	}
	startTime := time.Now()
	result, err := client.Invoke(&lambda.InvokeInput{
		FunctionName: &lambdaFunc,
		Payload:      payload})
	if err != nil {
		log.Printf("Error calling %s: %+v\n", lambdaFunc, err)
		return
	}
	elapsed := int(time.Since(startTime) / time.Millisecond)
	perf <- elapsed

	if debug {
		log.Printf("StatusCode: %d, elapsed time %d\n", result.StatusCode, elapsed)
		log.Printf("Returned orchestrator message: %s\n", string(result.Payload))
	}
}

func collectStats() {
	var count, aboveTarget, minMs, maxMs, total int
	startTime := time.Now()
	skip := 0
	elapsed := <-perf
	for elapsed >= 0 {
		if skip < 2*threads && elapsed > 500 {
			skip++
			elapsed = <-perf
			continue
		}
		count++
		total += elapsed
		if elapsed > targetMs {
			aboveTarget++
		}
		if elapsed > maxMs {
			maxMs = elapsed
		}
		if minMs == 0 || elapsed < minMs {
			minMs = elapsed
		}
		if count%printInterval == 0 {
			fmt.Printf("count %d slow %d min %d max %d avg %d current %d\n",
				count, aboveTarget, minMs, maxMs, total/count, elapsed)
		}
		elapsed = <-perf
	}
	fmt.Printf("count %d slow %d min %d max %d avg %d current %d\n",
		count, aboveTarget, minMs, maxMs, total/count, elapsed)
	fmt.Printf("Total elapsed time %s\n", time.Since(startTime))
}
