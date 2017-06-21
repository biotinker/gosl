#!/bin/bash

FILES="*.py"

echo
echo "monitoring:"
echo $FILES
echo
echo

while true; do
    inotifywait -q -e modify $FILES
    echo
    echo
    echo
    echo
    python zherk01.py
done
