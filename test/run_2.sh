#!/bin/sh

set -x

NICKNAME=green-sunset

go run ../ -config=../test/data/$NICKNAME.config.json