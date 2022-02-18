# TODO We must insert template-init into the env.
# - That means that tinit must be somewhere published.
# - I'll probably rename tinit to agent.
# - Once each environment is a separate VM the agent will run directly on the VM.

# FROM us-central1-docker.pkg.dev/devbookhq/devbook-runner-templates/base
FROM devbook

COPY files files

ENV root_dir="/home"
ENV start_cmd="echo 'Devbook'"

# ...
