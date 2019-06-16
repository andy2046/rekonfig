#!/usr/bin/env bash

set -euo pipefail

export GOCACHE=off

echo
echo === test pkg core ===
pkg="./pkg/core"
go test -v -race $pkg |\
    sed "s/PASS/$(printf "\033[32mPASS\033[0m")/"         |\
    sed "s/FAIL/$(printf "\033[31mFAIL\033[0m")/"         |\
    sed "s/RUN/$(printf "\033[33mRUN\033[0m")/"
go fmt $pkg
go vet $pkg
golint $pkg

echo
echo === test pkg controller ===
pkg="./pkg/controller/..."
go test -v -race $pkg |\
    sed "s/PASS/$(printf "\033[32mPASS\033[0m")/"         |\
    sed "s/FAIL/$(printf "\033[31mFAIL\033[0m")/"         |\
    sed "s/RUN/$(printf "\033[33mRUN\033[0m")/"
go fmt $pkg
go vet $pkg
golint $pkg
