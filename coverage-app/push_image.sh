#!/bin/bash
cd $( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

repo=coverage
image=${repo}:1.1

login=$(aws ecr get-login --no-include-email)
echo "login ECR"
( ${login} )

tag=${login##*/}

# check repository
aws ecr describe-repositories --repository-name ${repo} --query 'repositories[].repositoryUri' --output text
if [ $? -ne 0 ]; then
  echo "create repository ${repo} ..."
  aws ecr create-repository --repository-name ${repo}
fi

echo "push docker image to ECR ${tag}/${image}"
docker tag ${image} ${tag}/${image}
docker push ${tag}/${image}

sed -i -e "s|^IMAGE=.*|IMAGE=${tag}/${image}|" ./env.sh

# display result of ECR repo
aws ecr describe-images --repository-name ${repo}
