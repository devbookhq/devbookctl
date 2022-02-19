# ---- tinit build phase
FROM golang:1.17-alpine as build_go

COPY template-init /template-init

WORKDIR /template-init
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tinit .
# ---- tinit build phase END

# We need to base this image on the CUSTOMER custom image, so the ENV/ARG vars are defined and the tinit is executed as expected.
# Check what happens if the customer image does not defined run - does it inherit run from the old container? We may want to execute the run from VM - without defining it in container to prevent having to compose images.
# Check what happens with envs/args from the original image.
FROM ubuntu:18.04

ENV root_dir="/home/runner"
ENV start_cmd=";"

# Unix socket for communication with Runner.
# Used by the `tinit` process.
ENV runner_socket_path="/home/runner.socket"
# Go to the directory with setup files and execute the setup command - for examply installing dependencies, etc.

# Copy the startup script and set it to execute on container start.
COPY --from=build_go /tinit /usr/bin/tinit
CMD ["tinit"]
