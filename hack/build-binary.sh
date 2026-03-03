#!/bin/bash

set -euo pipefail
PROJECT_ROOT="$(realpath $(dirname $0)/..)"

HACK_DIR="$PROJECT_ROOT/hack"
source "$HACK_DIR/lib.sh"

VERSION=$("$HACK_DIR/get-version.sh")

echo "> Building binaries ..."
for pf in ${PLATFORMS//,/ }; do
  echo "> Building binary for $pf ..." | indent 1
  os=${pf%/*}
  arch=${pf#*/}
  CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -a -o "$PROJECT_ROOT/bin/kw-${os}.${arch}" \
    -ldflags "-X github.com/Diaphteiros/kw/pkg/version.StaticVersion=$VERSION" \
    "$PROJECT_ROOT" \
    | indent 1
done
