#!/bin/bash

cd "$( dirname "${BASH_SOURCE[0]}" )"
source env.sh

mycidr=$(curl ifconfig.me)/32

# set limited SSH rule for bastion host
bastionSgid=$(aws ec2 describe-security-groups --filters Name=group-name,Values=${EFS_STACK}-InstanceSecurityGroup* --query 'SecurityGroups[*].GroupId' --output text)
echo "set ssh rule for ${mycidr} in security group ${bastionSgid}"
aws ec2 revoke-security-group-ingress --group-id ${bastionSgid} --protocol tcp --port 22 --cidr 0.0.0.0/0
aws ec2 authorize-security-group-ingress --group-id ${bastionSgid} --protocol tcp --port 22 --cidr ${mycidr}

# set limited NFS rules for EFS mount
mountSgid=$(aws ec2 describe-security-groups --filters Name=group-name,Values=${EFS_STACK}-MountTargetSecurityGroup* --query 'SecurityGroups[*].GroupId' --output text)
nodeSgid=$(aws ec2 describe-security-groups --filters Name=group-name,Values=eksctl-${EKS_STACK}-nodegroup-0-SG* --query 'SecurityGroups[*].GroupId' --output text)
echo "set NFS rule for bastion sg ${bastionSgid} and node sg ${nodeSgid} in security group ${mountSgid}"
aws ec2 revoke-security-group-ingress --group-id ${mountSgid} --protocol tcp --port 2049 --cidr 0.0.0.0/0
aws ec2 authorize-security-group-ingress --group-id ${mountSgid} --protocol tcp --port 2049 --source-group ${bastionSgid}
aws ec2 authorize-security-group-ingress --group-id ${mountSgid} --protocol tcp --port 2049 --source-group ${nodeSgid}
