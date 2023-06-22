#!/bin/bash

# flash mega2560 with given ($1) .hex file

avrdude -v -c wiring -p m2560 -P /dev/ttyACM0 -b 115200 -D -F -U flash:w:$1