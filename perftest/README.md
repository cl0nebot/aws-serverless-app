# Benchmark Eligibility Request
This is a driver for testing the response time of eligibility orchestrator.

## Build and run benchmark test
```bash
make build
./test.sh
```
The test script uses `ssh` to the benchmark test on the bastion host, although you may also run the test from you local PC, e.g., for mac,
```bash
./perf-mac -debug
```
Or get more info on how to use the performance test driver.
```bash
./perf-mac -help
```

## Performance Test Result
![Performance](poc-perf.png)

This result shows that average response time of this service is 77-98 ms. The system scales linearly as the test driver increases message rate by using more threads.  Thus it can process at a rate of 500 message/s or more, as shown by the test case with 50 client threads.

The last 5 rows display sample response time of individual components, which shows that the org-status service is the major bottleneck of the overal system.  The orchestrator calls the first 3 services in parallel, i.e., Kafka publisher for request message, org-status service, and coverage service. The slowest of the 3 will be the main end-to-end delay of the end-to-end processing.

The org-status service is significantly slower than the coverage service because it calls 2 lambda functions sequencially, i.e., the Redis cache first and flogo rules second.  The response time of each lambda function call is not very predicable, it ranges from 25 to 100+ ms.  Thus each additional lambda function call may add up to 100 ms delay, although the actual process in the lambda function takes only less than 2 ms.

Although we may improve the performance of this POC further by combining the flogo-rules with the the flogo process for org-status, we kept these 2 processes as separate lambda functions to demostrate the overhead of each lambda function calls.