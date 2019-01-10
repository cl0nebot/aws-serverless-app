#!/bin/bash
SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source ${SDIR}/env.sh

curl "${GATEWAY_URL}" -d "{ \"resourceType\": \"EligibilityRequest\", \"ID\": \"1234\", \"patient\": {\"reference\": \"deceased\"},\"organization\": {\"reference\": \"P-000018\"},\"insurer\": {\"reference\": \"org2\"}, \"coverage\": {\"reference\": \"C-000022\"} }" 
echo ""
