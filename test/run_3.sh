#!/bin/sh

set -x

NICKNAME=black-dust

go run ../ -config=../test/data/$NICKNAME.config.json
