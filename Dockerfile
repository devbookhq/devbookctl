# ---- tinit build phase
FROM golang:1.17-alpine as build_go

COPY template-init /template-init

WORKDIR /template-init
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tinit .
#RUN go build -o /tinit .
# ---- tinit build phase END



# We need to base this image on the CUSTOMER custom image, so the ENV/ARG vars are defined and the tinit is executed as expected.
# Check what happens if the customer image does not defined run - does it inherit run from the old container? We may want to execute the run from VM - without defining it in container to prevent having to compose images.
# Check what happens with envs/args from the original image.
FROM ubuntu:18.04

RUN apt-get update && apt-get install -y \
  ca-certificates \
  curl \

# Define build arguments that are passed to `docker build` command.

# Relative path to template files in the Docker context.
ARG repo_files_dir

# Absolute path to directory where we will copy template files and then run `setup_cmd`.
ARG files_dir
ARG setup_cmd

# Absolute path to the directory where the template files will be copied when the container starts.
# It is also the path to the directory where you want to mount the host's volume and
# the place where you want to install new dependencies for the code cells.
ARG root_dir
# User defined command that will run in the `root_dir` after we copy files from `files_dir` and mount the host's volume there.
ARG start_cmd
# Subdirectory of the `root_dir` where we want to save the code cells' files.
ARG code_cells_dir

# Pass build arguments to env arguments so we can use them inside a running container.
# These variables are also accessible in the template's uploaded config file.
ENV files_dir=${files_dir}
ENV start_cmd=${start_cmd}
ENV root_dir=${root_dir}

# Unix socket for communication with Runner.
# Use by the `tinit` process.
ENV runner_socket_path="/home/runner.socket"

# Copy setup files from Docker context to the container.
COPY ${repo_files_dir} ${root_dir}

# Go to the directory with setup files and execute the setup command - for examply installing dependencies, etc.
WORKDIR ${root_dir}
RUN ${setup_cmd}

WORKDIR /

# Copy the startup script and set it to execute on container start.
COPY --from=build_go /tinit /usr/bin/tinit
CMD ["tinit"]
