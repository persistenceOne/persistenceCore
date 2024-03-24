#!/usr/bin/env bash

set -eo pipefail

echo "Generating gogo proto code"
cd proto

buf generate --template buf.gen.gogo.yaml $file

cd ..

# move proto files to the right places
mkdir protogen
cp -r github.com/* ./protogen/
rm -rf github.com
rm -rf protogen

go mod tidy
