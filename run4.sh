#!/bin/sh

set -x

NICKNAME=patient-haze

go run . -config=test/data/$NICKNAME.config.json
