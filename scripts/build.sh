#!/usr/bin/env bash

set -euo pipefail

echo
echo === build rekonfig image ===
operator-sdk build andy2046/rekonfig:latest
