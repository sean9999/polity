#!/bin/sh

set -x

NICKNAME=icy-cherry

cd ..

go run . -config=test/data/$NICKNAME.config.json
