#! /bin/bash

set -euxo pipefail

docker build \
--build-arg root_dir=/home/runner \
--build-arg start_cmd="" \
-t devbook .
