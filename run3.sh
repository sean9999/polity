#!/bin/sh

set -x

NICKNAME=icy-cherry

go run . -config=test/data/$NICKNAME.config.json
