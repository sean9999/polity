#!/bin/sh

set -x

NICKNAME=ice-cherry

go run . -config=test/data/$NICKNAME.config.json
