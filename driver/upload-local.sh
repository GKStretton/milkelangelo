#!/bin/bash
set -euo pipefail

SKETCH_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FQBN="esp32:esp32:esp32da"
PORT="${PORT:-/dev/ttyUSB0}"

LIBRARIES=(
  "AccelStepper@1.64.0"
)

for lib in "${LIBRARIES[@]}"; do
  arduino-cli lib install "$lib"
done

arduino-cli compile --fqbn "$FQBN" "$SKETCH_DIR"
arduino-cli upload --fqbn "$FQBN" --port "$PORT" "$SKETCH_DIR"
