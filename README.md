```
This project is under active development and is not ready for production use.
```

# Speedia AppManager

Speedia AppManager (SAM) is an open source application hosting manager in a single file. It has a REST API, CLI and dashboard built to help you run your applications in a container easily.

## Running

SAM is designed to manage an application with its dependencies in a lightweight container, specifically on top of [bitnami/minideb](https://hub.docker.com/r/bitnami/minideb/) image. However, it may also work as a regular binary on any Debian-based distro so you can use it to run your applications directly in a virtual machine or bare metal server.

To run SAM as a container, you can use the image available at DockerHub with the following command:

```
podman run --name sam --env 'VIRTUAL_HOST=speedia.net' --rm -p 10000:10000 -it speedia/sam:latest
```

Feel free to rename the container, vhost and change the host port as you wish. SAM should work with Docker, Docker Swarm, Rancher, Kubernetes, Portainer or any other tool that supports OCI-compliant containers.

You can publish port 80 and 443 to the host when running SAM in a virtual machine or bare metal server so that you don't need to use a reverse proxy, as long as your intention is to run a single application in the server.

Otherwise, you may want to use a reverse proxy to run multiple SAM instances in the same server and proxy each domain to the respective SAM instance, using [nginx-proxy/nginx-proxy](https://github.com/nginx-proxy/nginx-proxy) for example. Remember to also map port 10000 to a subdomain or directory in the reverse proxy for each SAM instance.

If you don't want to use containers, you can attempt to run SAM directly in a VM or server at your own risk. Download the latest release from the [releases page](https://github.com/speedianet/sam/releases) and use the [supervisord config file](https://github.com/speedianet/sam/blob/main/supervisord.conf) to run it as a service. In the future there will be a CLI command to automate this installation.

## Development

In this repository you'll find the REST API and CLI code plus the dashboard assets. The API and CLI uses Clean Architecture, DDD, TDD, CQRS, Object Calisthenics, etc. Understand how these concepts works before proceeding is advised.

To run this project during development you must install [Air](https://github.com/cosmtrek/air). Air is a tool that will watch for changes in the project and recompile it automatically.

### Environment Variables

You must have an `.env` file in the root of the git directory **during development**. You can use the `.env.example` file as a template. Air will read the `.env` file and use it to run the project during development.

If you add a new env var that is required to run the apis, please add it to the `src/presentation/shared/checkEnvs.go` file.

When running in production, the `/speedia/.env` file is only used if the environment variables weren't set in the system. For instance, if you want to set the `ENV1` variable, you can do it in the `.env` file or in the command line:

```
ENV1=XXX /speedia/sam
```

### Unit Testing

SAM commands can harm your system, so it's important to run the unit tests in a proper container:

```
podman build --format=docker -t sam-unit-test:latest -f Dockerfile.test .
podman run --rm -it sam-unit-test:latest
```

Make sure you have the `.env` file in the root of the git directory before running the tests.

### Dev Utils

The `src/devUtils` folder is not a Clean Architecture layer, it's there to help you during development. You can add any file you want there, but it's not recommended to add any file that is not related to development since the code there is meant to be ignored by the build process.

For instance there you'll find a `testHelpers.go` file that is used to read the `.env` during tests.

### Building

To build the project, run the command below. It takes two minutes to build the project at first. After that, it takes less than 10 seconds to build.

```
podman build --format=docker -t sam:latest .
```

To run the project you may use the following command:

```
podman run --name sam --env 'VIRTUAL_HOST=speedia.net' --rm -p 10000:10000 -it sam:latest
```

When testing, consider publishing port 80 and 443 to the host so that you don't need to use a reverse proxy.

## REST API

### Authentication

The API accepts two types of tokens and uses the standard "Authorization: Bearer \<token\>" header:

- **sessionToken**: is a JWT, used for dashboard access and generated with the user login credentials. The token contains the userId, IP address and expiration date. It expires in 3 hours and only the IP address used on the token generation is allowed to use it.

- **userApiKey**: is a token meant for M2M communication. The token is a _AES-256-CTR-Encrypted-Base64-Encoded_ string, but only the SHA3-256 hash of the key is stored in the server. The userId is retrieved during key decoding, thus you don't need to provide it. The token never expires, but the user can update it at any time.

### OpenApi // Swagger

To generate the swagger documentation, you must use the following command:

```
swag init -g src/presentation/api/api.go -o src/presentation/api/docs
```

The annotations are in the controller files. The reference file can be found [here](https://github.com/swaggo/swag#attribute).
