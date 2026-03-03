#!/bin/bash
# Requires git-extras (brew install git-extras)

set -euo pipefail
PROJECT_ROOT="$(realpath $(dirname $0)/..)"

HACK_DIR="$PROJECT_ROOT/hack"
source "$HACK_DIR/lib.sh"

VERSION=$("$HACK_DIR/get-version.sh")

echo "> Finding latest release"
major=${VERSION%%.*}
major=${major#v}
minor=${VERSION#*.}
minor=${minor%%.*}
patch=${VERSION##*.}
patch=${patch%%-*}
echo "v${major}.${minor}.${patch}"
echo

semver=${1:-"minor"}

case "$semver" in
  ("major")
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  ("minor")
    minor=$((minor + 1))
    patch=0
    ;;
  ("patch")
    patch=$((patch + 1))
    ;;
  (*)
    echo "invalid argument: $semver"
    exit 1
    ;;
esac

release_version="v$major.$minor.$patch"

echo "The release version will be $release_version. Please confirm with 'yes' or 'y':"
read confirm

if [[ "$confirm" != "yes" ]] && [[ "$confirm" != "y" ]]; then
  echo "Release not confirmed."
  exit 0
fi
echo

echo "> Updating version to release version"
"$HACK_DIR/set-version.sh" $major $minor $patch
echo

echo "> Creating release"
git release "$release_version"
# NESTED_MODULES must be set to the list of nested go modules, e.g. 'api,nested2,nested3'
for nm in ${NESTED_MODULES//,/ }; do
  echo "> Creating tag for $nm module"
  git tag "$nm/$release_version" -m "$nm/$release_version"
done
echo

echo "> Updating version to dev version"
dev_version="$release_version-dev"
"$HACK_DIR/set-version.sh" $major $minor $patch "dev"
echo

echo "> Pushing release"
git add --all
git commit -m "update version to $dev_version"
git push
for nm in ${NESTED_MODULES//,/ }; do
  echo "> Pushing tag for $nm module"
  git push origin tag "$nm/$release_version"
done
echo

echo "> Successfully finished"
