#!/bin/sh -e

export GOOS=${GOOS:-linux}
export GOARCH=${GOARCH:-arm64}
export VERSION="${VERSION:-v3.27.1}"

mkdir -p assets
cd assets
rm -rf goose
git clone --depth=1 -b ${VERSION} https://github.com/pressly/goose
cd goose
go mod tidy
go build -tags='no_mysql no_ydb no_postgres no_mssql no_vertica no_clickhouse no_libsql' -o goose.${GOOS} ./cmd/goose

cd ..
mv goose/goose.${GOOS} .
rm -rf goose

