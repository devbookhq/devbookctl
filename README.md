# devbookctl
devbookctl is a command-line tool for usedevbook.com


It allows you to build and push custom environments for Devbook VMs. You can then use VMs with your custom environments via the [Devbook SDK](https://github.com/devbookhq/sdk).

## Installation
TODO

## Usage

### Push custom VM environment
devbookctl expects a `Dockerfile.dbk` file to be present in the same directory from when you're calling the command.

The `Dockerfile.dbk` describe the VM's environments.

```sh
# 1. Go to a directory containing a Dockerfile describing your custom VM environment.

# 2. Call devbookctl.
$ devbookctl push
```

### Defining custom VM environmnet
See [example-env directory](./example-env).

Devbook VM environment is described via two files:
1. **`Dockerfile.dbk`**
- A Dockerfile describing the VM's environment. See more [here](#Dockerfile.dbk).
2. **`dbk.toml`**
- A configuration file. See more [here](#dbk.toml).

Both files must be present in the same directory from where you're calling the `devbookctl push` command.

#### `Dockerfile.dbk`
The Dockerfile describing the VM's environment.

The `Dockerfile.dbk` must start with `FROM devbook` declaration. This makes sure you use the base Docker image compatible with the Devbook VM.

```docker
# Use the Devbook base image.
FROM devbook

# Your custom setup.
# You can for example copy files, scripts, install packages, binaries, etc.
# E.g. install Python
RUN apt-get update && apt-get install -y \
  python3-pip
```

Note: Don't use the `CMD` or `ENTRYPOINT` commands in the Dockerfile. See section bellow on how to start a process as soon as a VM boots up.

#### `dbk.toml`
The [TOML](https://toml.io/en/) configuration file. The minimaln configuration file contains a `start_cmd` field.

```toml
# Required. Unique ID for your Devbook VM.
id = "example-env"

root_dir = "/home"

# Optional. A command that will get executed when the VM boots up.
# You can put your custom scripts here, for example.
start_cmd = "echo Hello World"
```

### Starting a process once the Devbook VM boots up
Use `start_cmd` in the `dbk.toml` configuration file to describe what command should be executed as soon as the VM boots up.
