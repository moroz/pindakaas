#!/usr/bin/env -S bash -euo pipefail

gitroot="$(git rev-parse --show-toplevel)"

export GOOS="${GOOS:-linux}"
export GOARCH="${GOARCH:-arm64}"

cd $gitroot
rm -rf rel
mkdir -p rel/assets

cd $gitroot/assets
pnpm run build

cd $gitroot
go build -o rel/server -tags PROD .
cp -R $gitroot/db/migrations rel/
cp $gitroot/db/seeds.sql rel/seeds.sql
install -m 0755 $gitroot/scripts/reset-staging-db.sh rel/reset-staging-db.sh
cp -R $gitroot/assets/dist rel/assets

TAR_OPTS="--no-xattrs"

if [[ "$(uname)" = "Darwin" ]]; then
  TAR_OPTS="--no-xattrs --no-mac-metadata"
fi

cd rel && tar czf release.tar.gz $TAR_OPTS server assets/ migrations/ seeds.sql reset-staging-db.sh
