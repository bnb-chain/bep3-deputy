#! /bin/bash

# Note: this script should be run from within its directory (otherwise it can't find the docker-compose.yaml file)

# exit on errors
set -e

# start up the deputy and chains
docker-compose up -d

# wait until the deputy is operational
while ! docker-compose exec deputy curl localhost:8080/status; do
 sleep 1
done

# run tests
# don't exit on error, just capture exit code (https://stackoverflow.com/questions/11231937/bash-ignoring-error-for-a-particular-command)
# use -count=1 to disable test result caching
go test . -count=1 -tags integration -v && exitStatus=$? || exitStatus=$?

# remove the deputy and chains
docker-compose down

exit $exitStatus
