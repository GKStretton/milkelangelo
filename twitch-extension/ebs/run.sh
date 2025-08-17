#!/bin/bash
set -euf -o pipefail

# The container running this script will send the script a SIGTERM when it shuts
# down. By default, this signal will not be passed down to child processes,
# and the container will continue to run until it eventually gets killed. We
# trap any appropriate signals and kill any child processes we started.
trap "trap - SIGTERM && pkill -e -TERM -ns 1" SIGINT SIGTERM EXIT

# Secrets that need setting
export SHARED_SECRET_GOO=${SHARED_SECRET_GOO}
export SHARED_SECRET_TWITCH=${SHARED_SECRET_TWITCH}

# deployed defaults
export BROADCAST_STATE_TO_TWITCH=${BROADCAST_STATE_TO_TWITCH:-true}
export ENABLE_SERVER_AUTHENTICATION=${ENABLE_SERVER_AUTHENTICATION:-true}

./ebs