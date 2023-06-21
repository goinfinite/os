```
This project is under active development and is not ready for production use.
```

# Speedia AppManager

Speedia AppManager (SAM) is an open source application hosting manager in a single file. It has a REST API, CLI and dashboard built to help you run your applications in a container easily.

In this repository you'll find the REST API and CLI code plus the dashboard assets. The API and CLI uses Clean Architecture, DDD, TDD, CQRS, Object Calisthenics, etc. Understand how these concepts works before proceeding is advised.

## Running

SAM is designed to manage an application with its dependencies in a container, specifically on top of Red Hat UBI 8. However, it may also work as a regular binary in any RPM-based distro so you can use it to run your applications directly in a virtual machine or bare metal server.

To run SAM as a container, you can use the image available at DockerHub with the following command:

```
podman run --name sam --env 'VIRTUAL_HOST=speedia.net' --rm -p 10000:10000 -it speedia/sam:latest
```

Feel free to rename the container, vhost and change the host port as you wish. SAM should work with Docker, Docker Swarm, Kubernetes, Portainer or any other tool that supports OCI-compliant containers.

You can publish port 80 and 443 to the host when running SAM in a virtual machine or bare metal server so that you don't need to use a reverse proxy, as long as your intention is to run a single application in the server.

Otherwise, you may want to use a reverse proxy to run multiple SAM instances in the same server and proxy each domain to the respective SAM instance, using [nginx-proxy/nginx-proxy](https://github.com/nginx-proxy/nginx-proxy) for example. Remember to also map port 10000 to a subdomain or directory in the reverse proxy for each SAM instance.

## Development

To run this project during development you must install [Air](https://github.com/cosmtrek/air). Air is a tool that will watch for changes in the project and recompile it automatically.

### Environment Variables

You must have an `.env` file in the root of the git directory **during development**. You can use the `.env.example` file as a template. Air will read the `.env` file and use it to run the project during development.

If you add a new env var that is required to run the REST API, please add it to the `src/presentation/api/checkEnvs.go` file.

When running in production, the `.env` file is ignored and the env vars must be set in the server/deployment or on the command line, for instance:

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

### Authentication Concepts

- **sessionToken**: is a temporary token that is used to identify a user. Mainly used in the dashboard, it is a JWT token that contains the user's account id, IP address and expiration date. It is usually stored in the browser's local storage.

- **userApiKey**: is a fixed token that is used to identify a user for API calls. It is usually stored in the user's environment variables.

```
Original String: U|1000|f5be9f20-1a26-44bc-87a5-33addccd4327
Encrypted String: GPDw4aeq+tDTTr987+xJLJdfjqC3Gm0QtSYnXYZp/X1ut2a9WpHMn9UnB0P8StWc+u+hunTyStvEWg=
```

The session token authentication is a standard JWT authentication, it's supposed to expire in 3 hours. Only the IP address used to create the token is allowed to use it.

The userApiKey is a _AES-256-CTR-Encrypted-Base64-Encoded_ string and is not stored after generating it. Only the SHA3-256 hash of the key is stored in file.

### OpenApi // Swagger

To generate the swagger documentation, you must use the following command:

```
swag init -g src/presentation/api/api.go -o src/presentation/api/docs
```

The annotations are in the controller files. The reference file can be found [here](https://github.com/swaggo/swag#attribute).
