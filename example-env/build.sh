#! /bin/bash

set -euxo pipefail

docker build \
-t example-env .
