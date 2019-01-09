#!/bin/bash
SDIR=$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source $(dirname "${SDIR}")/eks/aws/env.sh

cd ${SDIR}
echo "copy perftest driver to ${BASTION}"
scp -i ${SSH_PRIVKEY} -q -o "StrictHostKeyChecking no" ./perf ec2-user@${BASTION}:/home/ec2-user/

source ../orchestrator-app/env.sh

echo "run benchmark on ${BASTION}"
ssh -i ${SSH_PRIVKEY} -o "StrictHostKeyChecking no" ec2-user@${BASTION} << EOF
  ./perf -debug -count 10 -region ${AWS_REGION} -lambda ${FUNCTION_ARN}
EOF
