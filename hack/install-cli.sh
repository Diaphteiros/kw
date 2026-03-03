#!/bin/bash

set -euo pipefail
PROJECT_ROOT="$(realpath $(dirname $0)/..)"

HACK_DIR="$PROJECT_ROOT/hack"
source "$HACK_DIR/lib.sh"

VERSION=$("$HACK_DIR/get-version.sh")

echo "> Installing CLI ..."
CGO_ENABLED=0 go install \
  -ldflags "-X github.com/Diaphteiros/kw/pkg/version.StaticVersion=$VERSION" \
  "$PROJECT_ROOT" \
  | indent 1
