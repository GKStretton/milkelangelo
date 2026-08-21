#!/bin/bash
set -euo pipefail

SKETCH_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PORT="${PORT:-/dev/ttyUSB0}"

cd "$SKETCH_DIR"
pio run
pio run --target upload --upload-port "$PORT"
