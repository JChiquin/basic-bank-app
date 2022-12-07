#!/bin/bash

go list -json -m all | docker run --rm -i sonatypecommunity/nancy:latest sleuth --skip-update-check -q -o json > audit.json

VULNERABILITIES_COUNT=$(jq '.num_vulnerable' audit.json)

if [ $VULNERABILITIES_COUNT -gt 0 ]
then
    echo "Vulnerabilities found."

    jq -r '.vulnerable[].Coordinates' audit.json

    exit 1
fi

echo "Safe."
echo ""
