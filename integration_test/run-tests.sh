#! /bin/bash

# Note: this script should be run from within its directory (otherwise it can't find the docker-compose.yaml file)

# exit on errors
set -e

# start up the deputy and chains
docker-compose up -d

sleep 5
# run tests
# don't exit on error, just capture exit code (https://stackoverflow.com/questions/11231937/bash-ignoring-error-for-a-particular-command)
go test . -tags integration -v && exitStatus=$? || exitStatus=$?

# remove the deputy and chains
docker-compose down

exit $exitStatus
