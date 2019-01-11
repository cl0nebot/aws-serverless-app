# Benchmark Eligibility Request
This is a driver for testing the response time of eligibility orchestrator.

## Build and run benchmark test
```bash
make build
./test.sh
```
The test script uses `ssh` to run test on the bastion host, although you may also run the test from you local PC, e.g., for mac,
```bash
./perf-mac -debug
```
Or print out more info about how to use the performance test driver.
```bash
./perf-mac -help
```

## Performance Test Result
![Performance](poc-perf.png)

This result shows that average response time of this service is 77-98 ms. The system scales linearly as the test driver increases message rate by using more threads.  Thus it can process at a rate of 500 message/s or more, as shown by the test case with 50 client threads.

The last 5 rows in the above table display sample response time of individual components, which shows that the org-status service is the major bottleneck of the overall system.  The orchestrator calls the first 3 services in parallel, i.e., Kafka publisher for request message, org-status service, and coverage service. The slowest service of the 3 will determine the response time of the  end-to-end process.

The org-status service is significantly slower than the coverage service because it calls 2 lambda functions sequencially, i.e., the Redis cache first and then the flogo rules.  The response time of each lambda invocation is not very predictable, it ranges from 25 to 100+ ms.  Thus, each additional lambda function may add up to 100 ms delay, although the actual process in the lambda function takes only less than 2 ms.

Although we may improve the performance of this POC further by combining the flogo-rules with the flogo process for org-status, we kept these 2 processes as separate lambda functions to demostrate the overhead of each lambda invocation.

## AWS cost for development and test

Following is the AWS cost during the development and test of this POC.

![Cost](aws-cost.png)

Performance tests are done on Jan 7 and Jan 9 when we launched the following processes:
- 3 EC2 instances of type t2.medium in the EKS cluster
- 1 EC2 instance of type t2.medium as a bastion host
- 1 ElastiCache for Redis instance with node type cache.t2.small
- Deploy 3 Kafka brokers on the EKS cluster
- Deploy 3 containers for the coverage service on the EKS cluster
- Create 4 load-balancers for Kafka brokers (one for each broker, and one extra for external broker endpoint)
- Create 1 load-balancer as external endpoint for the coverage service

During the performance test days, the major AWS costs are
- EKS: This is the major expense charged when EKS cluster is up because we use CloudFormation stacks to create the cluster.  We could avoid this charge if we do not use CloudFormation, but instead, use KOPS to manage our own Kubernetes clusters.
- EC2-ELB: This is a major expense for the 5 ELB created for Kafka and the Coverage service.
- EC2-Other: Not sure what contributes to this charge.
- ElastiCache: This is charged whenever the Redis cache is up.  This POC uses only very minimal cache.  It would be much more expensive when the size and load of the cache increases for production.
- CloudWatch: This is charged for the logs of lambda functions.
- Others: This mainly includes EC2-Instances and CloudTrail. CloudTrail is a charge for AWS account management. The EC2-Instances would be much more expensive when we use more powerful EC2 instances for production.