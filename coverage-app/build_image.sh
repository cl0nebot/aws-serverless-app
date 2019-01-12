#!/bin/bash
SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
IMAGE=coverage:1.1

# verify BE_HOME
${BE_HOME}/bin/be-engine --help
if [ $? -ne 0 ]; then
  echo "Please set BE_HOME to TIBCO BE install root, e.g., /opt/tibco/be/5.5"
  exit 1
fi

cd ${SDIR}
if [ -f ./src/Coverage/coverage.cdd ]; then
  echo "copy coverage.cdd from src"
  cp ./src/Coverage/coverage.cdd ./build
fi

if [ ! -f ./build/Coverage.ear ]; then
  echo "build BE EAR file for coverage project"
  ${BE_HOME}/studio/bin/studio-tools --propFile ${BE_HOME}/studio/bin/studio-tools.tra -core buildEar -o ./build/Coverage.ear -p ./src/Coverage
fi

# generate Dockerfile if does not already exist
if [ ! -f ./build/Dockerfile ]; then
  if [ ! -f ./build/Coverage.ear ]; then
    echo "BE EAR file does not exist in ${SDIR}/build/Coverage.ear"
    exit 2
  fi
  if [ ! -f ./build/coverage.cdd ]; then
    echo "BE cdd file does not exist in ${SDIR}/build/coverage.cdd"
    exit 3
  fi
  echo "Generate Dockerfile for building image of Coverage service container ..."
  ${BE_HOME}/docker/bin/be-docker-gen --propFile ${BE_HOME}/docker/bin/be-docker-gen.tra -t ./build -i ubuntu-5.5.0-v01 -m POC -e poc@tibco.com -o true
  sed -i -e "s|^FROM .*|FROM shivchelwa/be-enterprise-base:ubuntu-5.5.0-v01|" ./build/Dockerfile
fi

# generate docker image for coverage
echo "Generate docker image coverage:1.0"
${BE_HOME}/docker/bin/build_app_image.sh ${IMAGE} ./build
