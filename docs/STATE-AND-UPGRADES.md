# State and Upgrades

Infinite OS is stateful by design. Applications, runtimes, databases, configuration, and files are installed after the environment starts and live on its filesystem. This page covers what that means for restarts, backups, and upgrades.

## Restarting and recreating

- `docker stop`, `start`, or `restart` keeps the same container and its filesystem;
- `docker rm` removes the container and its writable filesystem;
- creating another container from `docker.io/goinfinite/os:latest` creates a new environment; it does not inherit the old container's state.

## Upgrades

Pulling a newer Infinite OS image and recreating the container is a migration, not an in-place upgrade. Back up or snapshot the existing environment before you replace it.

[Infinite Ez](https://github.com/goinfinite/ez-issue-tracker/) can snapshot an entire Infinite OS environment for backup or migration. When you run Infinite OS standalone, use a container or host-level backup strategy until a dedicated standalone upgrade workflow is documented.

## Schedulers

Kubernetes and Swarm may recreate containers automatically. Do not deploy a stateful Infinite OS environment there without planning persistence first. An Infinite OS environment should normally run as a single active instance unless you have deliberately externalized and coordinated its state.
