# SpeediaOS AppManager API

SpeediaOS AppManager API is a REST API and CLI tool that allows you to manage your SpeediaOS and deployments.

This API uses Clean Architecture, DDD, TDD, CQRS, Object Calisthenics, etc. Understand how these concepts works before proceeding.

### Authentication Concepts

- **sessionToken**: is a temporary token that is used to identify a user. Mainly used in the dashboard, it is a JWT token that contains the user's account id, IP address and expiration date. It is usually stored in the browser's local storage.

- **userApiKey**: is a fixed token that is used to identify a user for API calls. It is usually stored in the user's environment variables.

```
Original String: U|1000|f5be9f20-1a26-44bc-87a5-33addccd4327
Encrypted String: GPDw4aeq+tDTTr987+xJLJdfjqC3Gm0QtSYnXYZp/X1ut2a9WpHMn9UnB0P8StWc+u+hunTyStvEWg=
```

The session token authentication is a standard JWT authentication, it's supposed to expire in 3 hours. Only the IP address used to create the token is allowed to use it.

The userApiKey is a _AES-256-CTR-Encrypted-Base64-Encoded_ string and is not stored after generating it. Only the SHA3-256 hash of the key is stored in file.

## Development

To run this project during development you must install [Air](https://github.com/cosmtrek/air). Air is a tool that will watch for changes in the project and recompile it automatically.

### Environment Variables

You must have an `.env` file in the root of the git directory **during development**. You can use the `.env.example` file as a template.

Air will read the `.env` file and use it to run the project during development.

If you add a new env var, please add it to the `src/presentation/api/checkEnvs.go` file. That file is used to check if all the required env vars were sent when running the binary in production.

When running in production, the `.env` file is ignored and the env vars must be set in the server/deployment or on the command line, for instance:

```
ENV1=XXX /speedia/sam-api
```

### Dev Utils

The `src/devUtils` folder is not a Clean Architecture layer, it's there to help you during development. You can add any file you want there, but it's not recommended to add any file that is not related to development since the code there is meant to be ignored by the build process.

For instance there you'll find a `testHelpers.go` file that is used to read the `.env` during tests.

### OpenApi // Swagger

To generate the swagger documentation, you must use the following command:

```
swag init -g src/presentation/api/main.go -o src/presentation/api/docs
```

The annotations are in the controller files. The reference file can be found [here](https://github.com/swaggo/swag#attribute).
