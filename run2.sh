#!/bin/sh

set -x

go run . -me="127.0.0.1:5001" -friend="127.0.0.1:5000"
