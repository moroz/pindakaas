#!/usr/bin/env -S bash -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-arm64}"

cd $gitroot
rm -rf rel
mkdir rel
cd $gitroot
go build -a -installsuffix cgo -o rel/server -tags PROD .
cp -R $gitroot/db/migrations rel/

# Remove hidden sql files, if any
rm $gitroot/db/migrations/._* || true

cd rel && tar czf release.tar.gz server migrations/

