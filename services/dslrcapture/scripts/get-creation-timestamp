#!/bin/bash
# returns the creation date and time of the file $1 in fractional seconds
# since epoch

# Get string timestamp, without timezone
str=$(stat -c %w $1 | sed "s/ +....//g")
# Turn this into [seconds].[nanos], printing to stdout
date -d "$str" '+%s.%N'
