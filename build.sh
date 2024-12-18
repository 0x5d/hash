#!/bin/bash

pushd src
ko build --platform linux/amd64 --local --preserve-import-paths --tag-only github.com/0x5d/hash
popd