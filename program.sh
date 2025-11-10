#!/bin/sh

set -e


(
    cd "$(dirname "$0")"
    go build -0 app/*.go
)

exec ./app "$@