#!/usr/bin/env bash

set -eo pipefail

go mod tidy

mkdir -p ./tmp_deps
mkdir -p ./tmp-swagger-gen

#copy some deps to use their proto files to generate swagger
deps="github.com/persistenceOne/pstake-native/v5 \
github.com/persistenceOne/persistence-sdk/v5 \
github.com/cosmos/cosmos-sdk \
github.com/cosmos/ibc-go/v10 \
github.com/CosmWasm/wasmd \
github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10 \
github.com/cosmos/ibc-apps/modules/ibc-hooks/v10"

set -- $deps
for dep in "$@"; do
  path=$(go list -f '{{ .Dir }}' -m $dep);
  cp -r $path tmp_deps;
done

chmod -R 755 tmp_deps

rm -rf tmp_deps/**/buf.work.yaml
rm -rf tmp_deps/**/testutil


proto_dirs=$(find ./proto ./tmp_deps -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  echo $dir
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf generate --template proto/buf.gen.swagger.yaml $query_file
  fi
done

swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp_deps
rm -rf ./tmp-swagger-gen
