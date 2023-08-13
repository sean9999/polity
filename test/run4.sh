#!/bin/sh

set -x

NICKNAME=patient-haze

cd ..

go run . -config=test/data/$NICKNAME.config.json
