#!/bin/bash

set -euo pipefail
PROJECT_ROOT="$(realpath $(dirname $0)/..)"

HACK_DIR="$PROJECT_ROOT/hack"
source "$HACK_DIR/lib.sh"

function tidy() {
  go mod tidy -e
}

# NESTED_MODULES must be set to the list of nested go modules, e.g. 'api,nested2,nested3'
for nm in ${NESTED_MODULES//,/ }; do
  echo "Tidy $nm module ..."
  (
    cd "$PROJECT_ROOT/$nm"
    tidy | indent 1
  )
done

echo "Tidy root module ..."
(
  cd "$PROJECT_ROOT"
  tidy | indent 1
)
