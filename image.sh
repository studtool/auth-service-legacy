#!/usr/bin/env bash

command="$1"
image="studtool/auth-service:$(head -n 1 version.txt)"

if [[ "$command" = "build" ]]; then
  docker build -t "$image" .
elif [[ "$command" = "push" ]]; then
  echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin \
    && docker push "$1"
elif [[ "$command" = "remove" ]]; then
  docker rmi "$image"
else
  echo "command expected 'build/push/remove'"
  exit -1
fi
