#!/bin/bash
# returns the creation date and time of the file $1 in fractional seconds
# since epoch

stat -c %w $1 | sed "s/ +....//g" | date '+%s.%N'

