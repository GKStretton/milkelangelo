#!/bin/bash
set -euf -o pipefail

# The container running this script will send the script a SIGTERM when it shuts
# down. By default, this signal will not be passed down to child processes,
# and the container will continue to run until it eventually gets killed. We
# trap any appropriate signals and kill any child processes we started.
trap "trap - SIGTERM && pkill -e -TERM -ns 1" SIGINT SIGTERM EXIT

export EBS_HOST=${EBS_HOST}
export SHARED_SECRET_EBS=${SHARED_SECRET_EBS}

export BROKER_HOST=${BROKER_HOST:-"localhost"}

./goo