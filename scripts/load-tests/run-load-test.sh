#!/bin/sh
# Source the envrc
echo "Pulling loadimpact docker image..."
docker pull loadimpact/k6  > /dev/null 2>&1

# Might need to add code to build image first
echo "Spinning up infura webserver"
# Perhaps put these server messages to logs
docker run -p 8000:8000 -d -t infura-web-server:latest 
# Let the server spin up
sleep 5

# Would like to run in docker container but 2 docker containers 
# don't like pinging eachother over localhost
# docker run -i loadimpact/k6 run - <${WORKSPACE}/scripts/load-tests/simple-health-check.js

k6 run ${1}

docker kill $(docker ps | grep infura-web-server | awk '{printf("%s", $1)}')
