```
This project is under active development and is not ready for production use.
```

# Speedia OS

Speedia OS is a container operating system designed to simplify the deployment and management of applications in a containerized environment. It offers a REST API, CLI, and a user-friendly dashboard for seamless container management.

## Running

To run Speedia OS as a container, you can pull the image from DockerHub and use the following command:

```
docker run --rm --name myapp-container \
  --env 'PRIMARY_VHOST=myapp.net' \
  -p 8080:80 -p 8443:443 -p 1618:1618 \
  -it docker.io/speedianet/os:latest
```

In this example, the container ports 80, 443, and 1618 are mapped to host ports 8080, 8443, and 1618, respectively. If you are running multiple containers on the same host, consider using a reverse proxy to manage traffic.

You can customize the container name, vhost, and host ports as needed. The --rm flag ensures the container is removed upon stopping. To retain the container, simply omit this flag.

After deploying the container, access the shell to create a new account with the following command:

```
docker exec -it myapp-container /bin/bash
os account create -u admin -p admin
```

Once the account is created, you can access the dashboard at `https://localhost:1618/_/` and log in with the credentials you just set up. Note that you may encounter an SSL warning due to the self-signed certificate, which you can ignore or replace with your own certificate later.

Through the dashboard, you can deploy applications using the Marketplace feature with just a few clicks. You can also use the CLI for deployments, such as:

```
os mktplace install -s wp -n myapp.net \
  -f 'adminUsername:admin' \
  -f 'adminPassword:abc123' \
  -f 'adminMailAddress:user@example.com'
```

The API Swagger documentation is available at `https://localhost:1618/_/api/swagger/`.

Speedia OS is compatible with Docker, Podman, Docker Swarm, Rancher, Kubernetes, Portainer, and any other tool that supports OCI-compliant containers.

## Development

In this repository you'll find the REST API and CLI code plus the dashboard assets. The API and CLI uses Clean Architecture, DDD, TDD, CQRS, Object Calisthenics, etc. Understand how these concepts works before proceeding is advised.

To run this project during development you must install [Air](https://github.com/cosmtrek/air). Air is a tool that will watch for changes in the project and recompile it automatically.

### Environment Variables

You must have an `.env` file in the root of the git directory **during development**. You can use the `.env.example` file as a template. Air will read the `.env` file and use it to run the project during development.

If you add a new env var that is required to run the apis, please add it to the `src/presentation/shared/checkEnvs.go` file.

When running in production, the `/speedia/.env` file is only used if the environment variables weren't set in the system. For instance, if you want to set the `ENV1` variable, you can do it in the `.env` file or in the command line:

```
ENV1=XXX /speedia/os
```

### Unit Testing

Speedia OS commands can harm your system, so it's important to run the unit tests in a proper container:

```
podman build -t sos-unit-test:latest -f Containerfile.test .
podman run --rm -it sos-unit-test:latest
```

Make sure you have the `.env` file in the root of the git directory before running the tests.

Some tests can run in your local machine, although it's not recommended. However, if you to give it a go, make sure to create the `/speedia/` directory before running the tests:

```
sudo mkdir /speedia
sudo chown $(whoami):$(whoami) /speedia
```

### Dev Utils

The `src/devUtils` folder is not a Clean Architecture layer, it's there to help you during development. You can add any file you want there, but it's not recommended to add any file that is not related to development since the code there is meant to be ignored by the build process.

For instance there you'll find a `testHelpers.go` file that is used to read the `.env` during tests.

### Building

To build the project, run the command below. It takes two minutes to build the project at first. After that, it takes less than 10 seconds to build.

```
podman build -t sos:latest .
```

To run the project you may use the following command:

```
podman run --name sos --env 'PRIMARY_VHOST=speedia.net' --rm -p 1618:1618 -it sos:latest
```

When testing, consider publishing port 80 and 443 to the host so that you don't need to use a reverse proxy.

### VSCode Extensions

The following extensions are highly encouraged to be used during development:

```
EditorConfig.EditorConfig
GitHub.copilot
GitHub.vscode-pull-request-github
esbenp.prettier-vscode
foxundermoon.shell-format
golang.go
hbenl.vscode-test-explorer
ms-vscode.test-adapter-converter
redhat.vscode-yaml
streetsidesoftware.code-spell-checker
streetsidesoftware.code-spell-checker-portuguese-brazilian
timonwong.shellcheck
```

## REST API

### Authentication

The API accepts two types of tokens and uses the standard "Authorization: Bearer \<token\>" header:

- **sessionToken**: is a JWT, used for dashboard access and generated with the account login credentials. The token contains the accountId, IP address and expiration date. It expires in 3 hours and only the IP address used on the token generation is allowed to use it.

- **accountApiKey**: is a token meant for M2M communication. The token is a _AES-256-CTR-Encrypted-Base64-Encoded_ string, but only the SHA3-256 hash of the key is stored in the server. The accountId is retrieved during key decoding, thus you don't need to provide it. The token never expires, but the user can update it at any time.

### OpenApi // Swagger

To generate the swagger documentation, you must use the following command:

```
swag init -g src/presentation/api/api.go -o src/presentation/api/docs
```

The annotations are in the controller files. The reference file can be found [here](https://github.com/swaggo/swag#attribute).
