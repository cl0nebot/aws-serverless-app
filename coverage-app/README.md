# coverage-app

This is a TIBCO BusinessEvents project to demostrate the deployment on Kubernetes. It implements coverage rules in either BE rules, or BE decision tables.  Decision table contains easy-to-edit rules in tabular format, and they can be hot-deployed at runtime.  

It also demonstrates the integration with Redis cache via either direct invocation of AWS lambda functions, or via REST API through the AWS API gateway. The direct lambda call is faster because it avoids the round-trip delay to the API gateway, which could save 100 ms or more.  To enable the direct lambda invocation, we have to configure the service role of the EKS cluster, which has been done in the script `eks/aws/efs-sg-rule.sh`.

In addition to Redis cache, this project also implemented rules to fetch coverage data from TIBCO ActiveSpaces, which is a fast distributed data grid for structured tuple data.  You may easily select one of the 3 mechanisms for data cache by editing the `CACHE_TYPE` in `coverage.yaml`. Valid values for this parameter are `REDIS_LAMBDA`, `REDIS_HTTP`, or `AS`.

The coverage service instances are deployed in muliple PODs, and the service is exposed by a LoadBalancer service.  Thus, the coverage service can be auto scaled as described in [kuberneties.io](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/).

## Develop BE applications

After you install the TIBCO BusinessEvents 5.5.0 or above, you can launch the BE studio, and import the project source code from `./src/Coverage`.  Then, right click the project root in the "Studio Explorer", and select "Properties". Select "Java Build Path", then edit the jars under the "Libraries" tab.  

Note that this project has compile-time dependency to 2 aws-java-sdk jars for direct invocation of AWS lambda functions.  You need to fix the path of these 2 jars to match your project location. One way to fix the path is to remove and then add them back with correct path.  This step can be automated if you use Maven.

If you are interested in how to setup Maven build and unit test for TIBCO BE projects, you may check the steps documented by [BE Developer's Guide](https://docs.tibco.com/pub/businessevents-enterprise/5.5.0/doc/pdf/TIB_businessevents-standard_5.5_developers_guide.pdf?id=7).  Or, for an alternative Maven setup, you may check the [be-sample](https://github.com/yxuco/be_sample/tree/master/SimpleHTTP) that contains a sample BE project and Maven POM files to build and test BE projects together with dependent Java projects and third-party jars.

To test this project with ActiveSpaces, you need to install TIBCO ActiveSpaces 3.5 and TIBCO FTL 5.4 or later, and also
* Edit `$BE_HOME/be-engine.tra` to set the following env:
```
tibco.env.FTL_HOME=/opt/tibco/ftl/5.4
tibco.env.AS3x_HOME=/opt/tibco/as/3.5
```
* Start AS data-grid and define a coverage table:
```
cd $AS_HOME/samples/scripts
../setup
./as-start
tibdg table create coverage orgid string
tibdg column create coverage effective_date string expire_date string
```

## Test BE application locally

This project uses a few third-party jars, which are under the `./build` folder.  You can build the project from the BE studio, or by using the script `${BE_HOME}/studio/bin/studio-tools`. An example for the use of the script is shown in `./build_image.sh`.  Put the result of the build, `Coverage.ear`, in the folder `./build`.

Set env `$BE_HOME` to your TIBCO installation folder, e.g., `/opt/tibco/be/5.5`. Edit the file `${BE_HOME}/bin/be-engine.tra` to add `./build` to the Java classpath, e.g.,
```bash
tibco.env.CUSTOM_EXT_PREPEND_CP=/path/to/poc/coverage-app/build
```
You can then start the BE engine locally by calling the script `./start-engine.sh`.

## Build and upload docker images
You can build the docker image for this BE application, and push it to AWS ECR using the following scripts, so the application can be deployed and started on EC2.
```bash
./build_image.sh
./push_images.sh
```

## Deploy and start coverage service
After the docker image is pushed to ECR, you can use the following script deploy and start the coverage service on the EKS cluster that we have already created using the script `eks/aws/create-all.sh`.
```bash
./deploy.sh
```
