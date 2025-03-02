#!/bin/bash
# upload via mqtt and pygateway

set -e

./scripts/build.sh
mosquitto_pub -h milkelangelo -t mega/flash -f ./build/firmware.ino.hex

echo "firmware sent to mega/flash"
