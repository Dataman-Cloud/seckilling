#!/usr/bin/env bash

DEFAULT_BUNDLES=(
binary
)

bundle() {
  local bundle="$1"; shift
  echo "---> Making bundle: $(basename "$bundle") (in $DEST)"
  source "script/$bundle" "$@"
}

if [ $# -lt 1 ]; then
   bundles=(${DEFAULT_BUNDLES[@]})
else
   bundles=($@)
fi
for bundle in ${bundles[@]}; do
  export DEST=.
  ABS_DEST="$(cd "$DEST" && pwd -P)"
  bundle "$bundle"
  echo
done

