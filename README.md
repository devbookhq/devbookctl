# `dbk`
`dbk` is a command-line tool (CLI) for [usedevbook.com](https://www.usedevbook.com/).

It allows you to build and push custom environments for Devbook VMs. You can then use VMs with your custom environments via the [Devbook SDK](https://github.com/devbookhq/sdk).


## Installation
Download `dbk` and install into `usr/local/bin`

```sh
curl -L https://usedevbook.com/install.sh | sh
```

## Usage - deploying custom environment
### Push a custom VM environment
To build and push a custom environment do the following:
```sh
# 1. Go to a directory containing a Dockerfile.dbk and dbk.toml describing your custom VM environment
$ cd <environment-directory>

# 2. Call dbk.
$ dbk push
```

### Defining a custom VM environment
See [example-env directory](./example-env).

Devbook VM environment is described via two files:
1. **`Dockerfile.dbk`** <br/>
A Dockerfile describing the VM's environment. See more [here](#Dockerfiledbk).
2. **`dbk.toml`** <br/>
A configuration file. See more [here](#dbktoml).

Both files must be present in the same directory from where you're calling the `dbk push` command.

#### `Dockerfile.dbk`
The Dockerfile describing the VM's environment.

The `Dockerfile.dbk` must start with `FROM devbook` declaration. This makes sure you use the base Docker image that is compatible with the Devbook's VM.

```dockerfile
# Use the Devbook base image.
FROM devbook

# Your custom setup.
# You can for example copy files, scripts, install packages, binaries, etc.
# E.g. install Python
RUN apt-get update && apt-get install -y \
  python3-pip
```

Note: Don't use the `CMD` or the `ENTRYPOINT` commands in the Dockerfile. See the section bellow on how to start a process as soon as a VM boots up.

#### `dbk.toml`
The [TOML](https://toml.io/en/) configuration file. The minimal configuration file contains just the `id` field.

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

## Usage - using a custom environment with the [Devbook SDK](https://github.com/devbookhq/sdk)
Pass the environment's `id` value from the `dbk.toml` config as the `env` parameter when initializing Devbook.

### React
```tsx
import { useDevbook } from '@devbookhq/sdk`

const { runCmd, stdout, stderr } = useDevbook({ env: 'example-env' })
```

### JavaScript/TypeScript
```ts
import { Devbook } from '@devbookhq/sdk`

const dbk = new Devbook({
  env: 'example-env`,
})
```
