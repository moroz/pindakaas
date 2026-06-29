#!/usr/bin/env -S bash -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-arm64}"

cd $gitroot
rm -rf rel
mkdir rel
cd $gitroot
go build -a -installsuffix cgo -o rel/server .
cp -R $gitroot/db/migrations rel/

TAR_OPTS="--no-xattrs"

if [[ "$(uname)" = "Darwin" ]]; then
  TAR_OPTS="--no-xattrs --no-mac-metadata"
fi

cd rel && tar czf release.tar.gz $TAR_OPTS server migrations/
