#! /bin/bash

set -euxo pipefail

# Try whether tinit is buildable
cd template-init && go build -o template-init/bin/tinit .
cd ..

docker build \
--build-arg repo_files_dir=templates/nextjs-v11-components/files \
--build-arg root_dir=/home/runner \
--build-arg start_cmd="npm run dev" \
--build-arg setup_cmd="npm i" \
-t template .

docker run \
-p 3000:3000 \
-v /tmp/runner.socket:/home/runner.socket \
-it template
# -it template /bin/ash