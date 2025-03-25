```
This project is still under active development (alpha stage).
Expect bugs early on. Create issues so they can be fixed.
```

# [Infinite OS](https://goinfinite.net/os/) &middot; [![Roadmap](https://img.shields.io/badge/roadmap-014737)](https://github.com/orgs/goinfinite/projects/9) [![Demo](https://img.shields.io/badge/read--only_demo-233876)](https://os.demo.goinfinite.net:1618/_/) [![/r/goinfinite](https://img.shields.io/badge/%2Fr%2Fgoinfinite-FF4500?logo=reddit&logoColor=ffffff)](https://www.reddit.com/r/goinfinite/) [![Discussions](https://img.shields.io/badge/discussions-751A3D?logo=github)](https://github.com/orgs/goinfinite/discussions) [![Report Card](https://img.shields.io/badge/report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/goinfinite/os) [![License](https://img.shields.io/badge/license-EPL-blue.svg)](https://github.com/goinfinite/os/blob/main/LICENSE.md)

Infinite OS is a container operating system designed to allow you to deploy applications knowing little to nothing about containers. It comes with a user-friendly dashboard, REST API, and CLI for seamless container management.

A read-only demo of the dashboard is available at [https://os.demo.goinfinite.net:1618/\_/](https://os.demo.goinfinite.net:1618/_/). The default credentials are `demo` and `abc123`.

## Running

To run Infinite OS you just need a single command:

```
docker run --rm --name myapp-container \
  --env 'PRIMARY_VHOST=myapp.net' \
  -p 8080:80 -p 8443:443 -p 1618:1618 \
  -it docker.io/goinfinite/os:latest
```

Then you'll be able to access the dashboard at `https://localhost:1618/_/` and the setup wizard will allow you to create a new account. Note that you may encounter an SSL warning due to the self-signed certificate, which you can ignore or replace with your own certificate later.

Using the dashboard you can deploy applications with the Marketplace feature in just a few clicks. You can also use the CLI for deployments, such as:

```
os mktplace install -s wp \
  -f 'adminUsername:admin' \
  -f 'adminPassword:abc123' \
  -f 'adminMailAddress:user@example.com'
```

In this example, the container ports 80, 443, and 1618 are mapped to host ports 8080, 8443, and 1618, respectively. If you are running multiple containers on the same host, consider using a reverse proxy to manage traffic or [Infinite Ez](https://github.com/goinfinite/ez), our free and easy-to-use self-hosted PaaS solution.

You can customize the container name, vhost, and host ports as needed. The `--rm` flag ensures the container is removed upon stopping. To retain the container, simply omit this flag.

The API Swagger documentation is available at `https://localhost:1618/api/swagger/`.

Infinite OS is compatible with Docker, Podman, Docker Swarm, Rancher, Kubernetes, Portainer, and any other tool that supports OCI-compliant containers.

## Development

The public roadmap for Infinite OS is available [here](https://github.com/orgs/goinfinite/projects/9). You may create issues or pull requests to contribute to the project.

In this repository you'll find the REST API and CLI code plus the dashboard assets. The API and CLI uses Clean Architecture, DDD, TDD, CQRS, Object Calisthenics, etc. Understand how these concepts works before proceeding is advised.

### Building

#### Simple Build

To build the project, run the command below. It takes two minutes to build the project at first. After that, it takes less than 10 seconds to build.

```
podman build -t os:latest .
```

To run the project you may use the following command:

```
podman run --name os --env 'PRIMARY_VHOST=goinfinite.local' --rm -p 1618:1618 -it os:latest
```

When testing, consider publishing port 80 and 443 to the host so that you don't need to use a reverse proxy. You should also consider using `--env 'LOG_LEVEL=debug'` to increase the log verbosity.

#### Development Build

When developing the project, you may want to use a script to automate the build process. The `dev-build.sh` script is available in the root of the project and it will take care of all the steps needed to build and run the container.

To run the script you can simply use `bash dev-build.sh` (bash may be replaced by zsh or similar). By default, the script will expose the port 1618 to the host which is used by the API and the dashboard.

- If you pass the `http` argument, it will also expose the ports 80 and 443 to the host;
- If you pass the `ols` argument, it will expose port 7080 (used by OpenLiteSpeed admin);
- If you pass the `no-cache` argument, it will remove the image cache and rebuild the image from scratch;

The script will also create a `dev` account with the password `123456` so you can access the dashboard.

When you need to stop the container, just CTRL+C to stop and remove it. If you don't want to remove it, just ditch the `--rm` flag from the `podman run` command in the script.

If you look closely at the script, you'll see that it mounts the project's `bin` directory to the container `/infinite/bin` path. This is done to allow the container to access the binary file generated by Air on the host. The script then replace the binary file that comes with the container with the one on the `/bin` directory.

With this approach you don't need to rebuild the container every time you change the code. Although sometimes you may want to restart the container to apply some changes, specially when changing the dependencies or system configurations. In this case, just hit CTRL+C to stop the container and run the script again.

**Notes:**

1. You must run the script from the project's root directory;
2. Until Echo v4.13.0 is released, you'll need to refresh the browser page during development to see the changes in the dashboard as we're not able to use the `DEV_MODE` auto refresh websocket trick for now. To understand how this trick used to work, check the UI router and main layout files.

### Unit Testing

Infinite OS commands can harm your system, so it's important to run the unit tests in a proper container:

```
podman build -t os-unit-test:latest -f Containerfile.test .
podman run --rm -it os-unit-test:latest
```

Make sure you have a `.env` file in the root of the git directory before running the tests.

Some tests can run in your local machine, although it's not recommended. However, if you to give it a go, make sure to create the `/infinite/` directory before running the tests:

```
sudo mkdir /infinite
sudo chown $(whoami):$(whoami) /infinite
```

### Dev Utils

The `src/devUtils` folder is not a Clean Architecture layer, it's there to help you during development. You can add any file you want there, but it's not recommended to add any file that is not related to development since the code there is meant to be ignored by the build process.

For instance there you'll find a `testHelpers.go` file that is used to read the `.env` during tests.

### Web UIs

This project has two web UIs, the previous Vue.js frontend and the new [Templ](https://templ.guide/) + [Alpine.js](https://alpinejs.dev/) + [HTMX](https://htmx.org/docs/) frontend. The Vue.js frontend is deprecated and will be removed in the future. It's available at `/_/` and the [Templ](https://templ.guide/) + [Alpine.js](https://alpinejs.dev/) + [HTMX](https://htmx.org/docs/) frontend is available at `/`.

The new frontend based on the [Templ](https://templ.guide/) + [Alpine.js](https://alpinejs.dev/) + [HTMX](https://htmx.org/docs/) combo mentioned was developed as a proof of concept to create an interface without needing to leave Go. To understand the entire conceptual and theoretical foundation behind using these technologies to create a new architecture, [access this article](https://ntorga.com/full-stack-go-app-with-htmx-and-alpinejs/). However, to grasp the practical basis of how to apply this new architecture, [refer to the proof of concept](https://github.com/ntorga/clean-ddd-full-stack-go-poc) used to develop it.

For the interface code to be read and rendered by Go, we need to convert all `.templ` files into `.go` files. To do this, run the following command at the root of the application:

```
templ generate -path src/presentation/api
```

It is important that this is done before using Air to create the binary; otherwise, the Web UI will not be embedded, and you will not be able to use it.

**NOTE:** If you are using the `dev-build.sh` script, you don't need to run the `templ generate` command (or Air for that matter) since the script will take care of everything for you.

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
timonwong.shellcheck
```

## REST API

### Authentication

The API accepts two types of tokens and uses the standard "Authorization: Bearer \<token\>" header:

- **sessionToken**: is a JWT, used for dashboard access and generated with the account login credentials. The token contains the accountId, IP address and expiration date. It expires in 3 hours and only the IP address used on the token generation is allowed to use it.

- **accountApiKey**: is a token meant for M2M communication. The token is an _AES-256-CTR-Encrypted-Base64-Encoded_ string, but only the SHA3-256 hash of the key is stored in the server. The accountId is retrieved during key decoding, thus you don't need to provide it. The token never expires, but the user can update it at any time.

### OpenApi // Swagger

To generate the swagger documentation, you must use the following command:

```
swag init -g src/presentation/api/api.go -o src/presentation/api/docs
```

The annotations are in the controller files. The reference file can be found [here](https://github.com/swaggo/swag#attribute).
