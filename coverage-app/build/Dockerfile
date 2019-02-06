#
# Dockerfile for project: Coverage.ear
#

FROM shivchelwa/be-enterprise-base:ubuntu-5.5.0-v01
MAINTAINER POC poc@tibco.com

ARG BE_APP_NAME
ARG BE_APP_VERSION
ARG TARGET_DIR

LABEL BE_APP_NAME=$BE_APP_NAME BE_APP_VERSION=$BE_APP_VERSION

ADD . /opt/tibco/be/application

RUN mkdir -p /opt/tibco/be/application/ear
ADD Coverage.ear /opt/tibco/be/application/ear
RUN rm /opt/tibco/be/application/Coverage.ear

ENV CDD_FILE /opt/tibco/be/application/coverage.cdd
ENV EAR_FILE /opt/tibco/be/application/ear/Coverage.ear

#External Jars
ADD httpcore-4.4.9.jar /opt/tibco/be/ext
ADD httpclient-4.5.5.jar /opt/tibco/be/ext
ADD jackson-core-2.6.7.jar /opt/tibco/be/ext
ADD jackson-databind-2.6.7.2.jar /opt/tibco/be/ext
ADD aws-java-sdk-core-1.11.478.jar /opt/tibco/be/ext
ADD joda-time-2.8.1.jar /opt/tibco/be/ext
ADD jackson-annotations-2.6.0.jar /opt/tibco/be/ext
ADD aws-java-sdk-lambda-1.11.478.jar /opt/tibco/be/ext

#Application's environment variables
ENV AWS_REGION us-west-2
ENV Deployment Coverage
ENV HTTP_PORT 8890
ENV WAYPOINTS_TOPIC waypoints
ENV AS_URL http://localhost:8080
ENV Domain domain
ENV MessageEncoding ISO8859-1
ENV REDIS_LAMBDA_ARN arn:aws:lambda:us-west-2:742759186184:function:coverage-reference-app-CoveragedataFunction-SUM238IKHALO
ENV REDIS_HTTP_URL https://es0r4xjoeb.execute-api.us-west-2.amazonaws.com/Prod/coveragedata/
ENV KAFKA_URL af9d687aa120311e9b3520a2b998cd74-1731976297.us-west-2.elb.amazonaws.com:9094
ENV CACHE_TYPE REDIS_LAMBDA

#Exposing ports
EXPOSE 8890 
