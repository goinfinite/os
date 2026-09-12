---
name: os-usage
description: Use when deploying or managing applications on an Infinite OS instance — maps the dashboard, CLI, and REST API, and gives the core deployment workflows.
version: 1.0.0
lastUpdated: 2026-09-11
---

## Purpose

Infinite OS is a wildcard container image. It gives you a host-in-a-box experience: applications, the services they need, databases, a web server, and a scheduler, all managed from one dashboard, CLI, and REST API. The container becomes the host you need after you run it.

Use this skill to deploy an instance, deploy WordPress or a custom application, install services, map hostnames, issue TLS certificates, create databases, and schedule jobs.

## Procedure

### 1. Check your position

Before you deploy anything, find out where you are and what is available:

1. Run `os version`. If it answers, you are inside an Infinite OS container. Skip to step 3.
2. Run `command -v docker` and `command -v podman`. If neither answers, stop and ask the user to install Docker or Podman, or to run you inside an Infinite OS container. Do not install a container runtime on your own.
3. List the running containers with `docker ps` or `podman ps`. If an Infinite OS container already runs, note its name (every command below will use it as the `docker exec <container>` prefix) and skip to step 3.
4. Otherwise, deploy an instance in step 2.

### 2. Deploy Infinite OS

```sh
docker run -d --name 'myapp.net' \
  --env 'PRIMARY_VHOST=myapp.net' \
  -p 8080:80 -p 8443:443 -p 1618:1618 \
  docker.io/goinfinite/os:latest
```

- `PRIMARY_VHOST` sets the domain the instance serves. The container auto-detects it from the hostname if unset.
- Port `1618` serves the dashboard and the API. It is the only required port.
- Ports `80` and `443` serve the deployed applications directly. Publish them when you do not use a reverse proxy. Map them to `8080` and `8443` on the host when the host ports are taken.
- Run the container as a long-lived instance. The container itself is the unit of state: the whole filesystem holds the OS databases, the deployed applications, and their data. Do not treat the container as disposable.
- The container name usually matches the domain.
- Any OCI runtime works: Docker, Podman, Kubernetes, and others. Container platforms such as Infinite Ez, Kubernetes, and Coolify can run the image too. Infinite Ez, our PaaS, runs one OS container per account, like cPanel runs one account per customer, and can snapshot the entire container to a downloadable `.zip` to move it to another instance.
- Replace `-d` with `-it` to run an attached session.

Open `https://localhost:1618/`. The first login runs a setup wizard that creates the admin account. The instance uses a self-signed certificate, so the browser warns once.

### 3. Run OS commands

The CLI is the recommended path for an agent. It needs no authentication. It runs as root, and `os` refuses to run as any other user.

> **Outside the container, prefix every command with the container name.** Run `docker exec <container> os <command>` instead of `os <command>`, for example `docker exec myapp.net os services get`. Copy files with `docker cp`. Inside the container, run the commands as written.

- `os <command> --help` lists the flags and marks the required ones.
- `os mktplace list-catalog` lists deployable applications with their slug, id, and data fields. The list paginates at 50 items per page; filter with `-s <slug>` or raise `--items-per-page`.
- `os services get-installables` lists installable services: runtimes, databases, and web servers. The list paginates at 10 items per page; raise `--items-per-page` to see them all.

Use the API only when the container shell is out of reach. Use the dashboard for a human.

### 4. Deploy a marketplace application

```sh
os mktplace install -s wp -n myapp.net \
  -f 'adminUsername:admin' -f 'adminPassword:...' \
  -f 'adminMailAddress:user@example.com'
```

- `-s` is the catalog slug, for example `wp` for WordPress, and `-i` the catalog id. Use one of them. List both with `os mktplace list-catalog`.
- `-n` sets the hostname. Omit it to use the primary virtual host.
- `-d` sets the URL path. Omit it to serve at the hostname root.
- `-f key:value` passes one installation data field. Repeat it. The catalog lists each item's `dataFields`, with `name`, `isRequired`, `defaultValue`, and `options`.

The install pulls the dependencies, creates the service, starts it, and creates the mapping. Check `os mktplace list` and `os services get`. Remove an installation with `os mktplace delete -i <installedId>`.

### 5. Deploy a custom application

Clone or copy the project into the container. The container ships `git`, so clone directly:

```sh
docker exec <container> git clone <repoUrl> /app/myapp
```

Every hostname other than the primary virtual host must exist before a mapping or a service `-H` references it. Create it with `os vhost create -n <hostname> -t top-level`.

#### PHP

Install the php-webserver runtime, map the hostname to it, and pick the PHP version. The mapping creates the PHP virtual host, which serves the vhost root. Keep the project in `/app/myapp` and replace the vhost directory with a symlink to its public directory:

```sh
os vhost create -n blog.myapp.net -t top-level
docker exec <container> rmdir /app/html/blog.myapp.net
docker exec <container> ln -s /app/myapp/public /app/html/blog.myapp.net
os services create-installable -n php-webserver -v 8.3
os vhost mapping create -n blog.myapp.net -p / -t service -v php-webserver
os runtime php update -n blog.myapp.net -v 8.3
```

A flat PHP app can live directly in `/app/html/<hostname>/` instead. For a Composer project, install Composer after php-webserver:

```sh
docker exec <container> php -r "copy('https://getcomposer.org/installer', '/tmp/composer-setup.php');"
docker exec <container> php /tmp/composer-setup.php --install-dir=/usr/local/bin --filename=composer
docker exec <container> rm /tmp/composer-setup.php
docker exec <container> bash -c 'cd /app/myapp && composer install'
```

#### Laravel

Create the database with `os db create` and `os db create-user`, then prepare the app and make its writable directories writable by the web server (`nobody:nogroup`):

```sh
docker exec <container> bash -c 'cd /app/myapp && cp .env.example .env && php artisan key:generate'
docker exec <container> chown -R nobody:nogroup /app/myapp/storage /app/myapp/bootstrap/cache
```

Set `DB_HOST=127.0.0.1`, `DB_PORT=3306`, `DB_DATABASE`, `DB_USERNAME`, and `DB_PASSWORD` in `.env` to match the database, then run the migrations:

```sh
docker exec <container> bash -c 'cd /app/myapp && php artisan migrate --force'
```

Enable any PHP module the app needs with `os runtime php update-modules`.

#### Node.js, Bun, Python, Ruby, and Java

The service runs a startup file on a fixed port. Create the virtual host, then install the runtime, point it at the project entry file, and set the working directory:

```sh
os vhost create -n blog.myapp.net -t top-level
os services create-installable -n node -v lts -f /app/myapp/server.js -w /app/myapp -e PORT=3000 -H blog.myapp.net
```

The app must listen on the runtime port: `3000` for Node, Bun, and Ruby, `8000` for Python, `8080` for Java. `lts` in `-v` is an alias for the current LTS release; `os services get-installables` lists concrete versions too. A service `-p` is the internal port nginx proxies to, not a host-published port.

Pass the port and other environment variables with `-e name=value`, for example `-e PORT=3000`. For a WebSocket app, register the port as `ws`, for example `-p '3000/ws'`, so nginx forwards the upgrade headers. When the app needs its own server command, override the start command with `-c`, for example FastAPI:

```sh
os services create-installable -n python -v 3.12 -f /app/myapp/main.py \
  -c '/app/myapp/.venv/bin/uvicorn main:app --host 0.0.0.0 --port 8000' \
  -w /app/myapp -H blog.myapp.net
```

Install the app dependencies after the runtime, then restart the service (find its name with `os services get`). Node uses `npm install` and Ruby `bundle install`:

```sh
docker exec <container> bash -c 'cd /app/myapp && mise x node@lts -- npm install'
os services update -n <serviceName> -s restart
```

Python uses `uv`. It is not part of the runtime, so install it with mise, create a project virtual environment, and install the dependencies into it:

```sh
docker exec <container> bash -c 'cd /app/myapp && mise install uv && mise x python@3.12 uv -- uv venv && mise x uv -- uv pip install -r requirements.txt'
```

Start the app from the virtual environment, for example `-c '/app/myapp/.venv/bin/uvicorn main:app --host 0.0.0.0 --port 8000'`.

For Rails, run the migrations with the Ruby toolchain and point `config/database.yml` at `127.0.0.1:5432`:

```sh
docker exec <container> bash -c 'cd /app/myapp && mise x ruby@3.4 -- bin/rails db:migrate'
```

#### Django

Run it with gunicorn on the runtime port. Add gunicorn to `requirements.txt`, then override the start command and install the dependencies:

```sh
os vhost create -n blog.myapp.net -t top-level
docker exec <container> bash -c 'cd /app/myapp && mise install uv && mise x python@3.12 uv -- uv venv && mise x uv -- uv pip install -r requirements.txt'
os services create-installable -n python -v 3.12 \
  -c '/app/myapp/.venv/bin/gunicorn myproject.wsgi:application --bind 0.0.0.0:8000' \
  -w /app/myapp -H blog.myapp.net
docker exec <container> bash -c 'cd /app/myapp && /app/myapp/.venv/bin/python manage.py migrate'
```

Add the hostname to `ALLOWED_HOSTS` and point the database settings at `127.0.0.1:5432`. Set `STATIC_ROOT` to `staticfiles`, collect the static files, and serve them through a static mapping:

```sh
docker exec <container> bash -c 'cd /app/myapp && /app/myapp/.venv/bin/python manage.py collectstatic --noinput'
docker exec <container> ln -s /app/myapp/staticfiles /app/html/blog.myapp.net/static
os vhost mapping create -n blog.myapp.net -p /static -t static-files
```

Gunicorn is the production default for synchronous Django. Async/ASGI Django uses Uvicorn (often as Gunicorn workers) or Daphne. Granian is the newer Rust-based server that supports both interfaces.

#### Static sites

Create the virtual host and put the files in its public directory. The mapping falls back to `index.html` for unknown paths, so client-side routing works:

```sh
os vhost create -n blog.myapp.net -t top-level
docker exec <container> git clone <repoUrl> /app/html/blog.myapp.net
os vhost mapping create -n blog.myapp.net -p / -t static-files
```

For a site that needs a build, clone it outside the vhost root, build it, and copy the output in:

```sh
docker exec <container> git clone <repoUrl> /app/myapp
docker exec <container> bash -c 'cd /app/myapp && mise x node@lts -- npm install && mise x node@lts -- npm run build'
docker exec <container> bash -c 'cp -r /app/myapp/dist/. /app/html/blog.myapp.net/'
```

This covers React, Vue, Svelte, Quasar, and every other frontend that builds to static files. The output directory varies: `dist/` for Vite-based builds, `dist/spa/` for Quasar, `build/` for Create React App, `out/` for a Next.js static export.

Next.js and Nuxt render on the server by default, so they run as Node services instead. Create the service with the framework start command, install the dependencies, build, and restart:

```sh
os services create-installable -n node -v lts -c 'mise x node@lts -- npm start' -w /app/myapp -e PORT=3000 -H blog.myapp.net
docker exec <container> bash -c 'cd /app/myapp && mise x node@lts -- npm install && mise x node@lts -- npm run build'
os services update -n <serviceName> -s restart
```

For Next.js, `npm start` runs `next start`. For Nuxt, start `node .output/server/index.mjs` after the build, and set the cache URL with `-e REDIS_URL=redis://127.0.0.1:6379`.

This path uses nginx and ignores `.htaccess`. If the site needs `.htaccess`, use the php-webserver runtime instead. It serves static files and honors `.htaccess` rewrites:

```sh
os services create-installable -n php-webserver -v 8.3
os vhost mapping create -n blog.myapp.net -p / -t service -v php-webserver
```

#### Compiled applications

Build the artifact first. Install a catalog runtime with `os services create-installable` and use its toolchain; `mise install` is only for a toolchain the catalog lacks, such as Go or Rust. You can also build on the host and copy the result in:

```sh
docker exec <container> bash -c 'cd /app/myapp && mise install go@1.27.1 && mise x go@1.27.1 -- go build -o app .'
docker exec <container> bash -c 'cd /app/myapp && mise install rust@latest && mise x rust@latest -- cargo build --release'
```

The Go command builds the main package at the repo root; adjust the path for layouts like `./cmd/app`. The Rust binary lands at `target/release/<binary>`. A Java JAR needs the Java runtime, then runs with `mise x java@21 -- java -jar /app/myapp/app.jar`:

```sh
os services create-installable -n java -v 21
```

Then create the virtual host and a custom service that runs the binary:

```sh
os vhost create -n blog.myapp.net -t top-level
os services create-custom -n myapp -t other -c '/app/myapp/app' \
  -w /app/myapp -p '8080/http' -a true -H blog.myapp.net
```

- `-t` is the service type: `database`, `runtime`, `webserver`, or `other`.
- `-c` is the start command and `-w` the working directory.
- `-p port/protocol` registers the port the service listens on.
- `-e name=value` sets an environment variable. Repeat it.
- `-a true` creates the mapping; `-H` sets its hostname and defaults to the primary virtual host at `/`.
- When the start command needs an interpreted runtime, wrap it in the toolchain runner, for example `-c 'mise x node@lts -- npm start'`.

For a background worker with no web interface, pass `-a false` and no `-p` or `-H`. The service still starts and restarts on failure, and writes its log to `/app/logs/<serviceName>/<serviceName>.log`.

Check the service status with `os services get` and the mapping with `os vhost mapping get`. Then request the hostname to confirm the app answers, for example `curl -k https://blog.myapp.net/health`.

### 6. Expose the application

The marketplace install and the runtime install already create the mapping. `create-custom` does too when `-a true`. Create a mapping by hand only for extra paths:

```sh
os vhost mapping create -n myapp.net -p /blog -t service -v <serviceName>
```

- `-t` is the target type: `url`, `service`, `response-code`, `inline-html`, or `static-files`.
- `-v` is the target value, for example the service name or a URL.
- `-m` sets the match pattern: `begins-with`, `contains`, `equals`, or `ends-with`.

List mappings with `os vhost mapping get`. Issue a publicly trusted certificate with:

```sh
os ssl create-trusted -n myapp.net
```

The hostname must resolve to the host and be reachable on port 80. For a certificate pair you already have, use `os ssl create -v myapp.net -c <certPath> -k <keyPath>`. List certificates with `os ssl get`.

### 7. Support the application

```sh
os db create -t mariadb -n appdb
os db create-user -t mariadb -n appdb -u appuser -p '...' -r 'ALL'
os cron create -s '0 3 * * *' -c '<command>' -d 'nightly backup'
```

`-t` accepts `mariadb` (alias `mysql`) and `postgresql` (alias `postgres`). Install MongoDB, Redis, and OpenSearch as services with `os services create-installable`, for example `os services create-installable -n redis -v 7.2`. Apps connect to `127.0.0.1` on the engine port: `3306` MariaDB, `5432` PostgreSQL, `27017` MongoDB, `6379` Redis. The cron `-s` field uses the standard five-field syntax. List databases with `os db get` and jobs with `os cron get`.

Cron runs the command as root with a minimal environment and no retry. Use absolute paths and redirect the output, for example `-c '/app/myapp/.venv/bin/python /app/myapp/pipeline.py >> /app/logs/cron/pipeline.log 2>&1'`.

### 8. Tune the PHP runtime

```sh
os runtime php update -n myapp.net -v 8.3
os runtime php update-setting -n myapp.net -v 8.3 -N memory_limit -V 512M
os runtime php update-modules -n myapp.net -v 8.3 -m '<module>:true'
```

Read the current configuration with `os runtime php get -n <hostname>`.

### 9. Manage files and accounts

The application files live on the container filesystem. Use `docker cp` or a shell. The dashboard at `/file-manager/` and the API at `/api/v1/files/` offer the same operations without shell access: `GET /` reads, `POST /` creates, `PUT /` updates, `POST /upload/`, `GET /download/`, `POST /compress/`, `PUT /extract/`, `PUT /delete/`.

Accounts and SSH access keys use `os account get|create|update|delete|create-public-key|delete-key`.

### 10. Use the dashboard and the API

- Dashboard: `https://<host>:1618/`. Pages: `/marketplace/`, `/runtimes/`, `/databases/`, `/mappings/`, `/ssls/`, `/crons/`, `/file-manager/`, `/accounts/`, `/overview/`.
- REST API: `https://<host>:1618/api/v1/...`.

Get a token:

```sh
curl -k -X POST 'https://<host>:1618/api/v1/auth/login/' \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"..."}'
```

The token is in the response body at `body.tokenStr`. Send it on every later call:

```sh
curl -k 'https://<host>:1618/api/v1/...' -H "Authorization: Bearer $TOKEN"
```

The instance uses a self-signed certificate by default, so pass `-k`. Every controller response has the same envelope: `status`, `body`, `readableMessage`. The `status` is the HTTP code, so a `2xx` status means the call worked. Middleware and parse errors can use a different shape. The Swagger UI at `https://<host>:1618/api/swagger/index.html` documents every route and body.

## Interface map

- Deploy a marketplace app: `os mktplace install` / `POST /api/v1/marketplace/catalog/`
- Install a service: `os services create-installable` / `POST /api/v1/services/installables/`
- Create a custom service: `os services create-custom` / `POST /api/v1/services/custom/`
- Create a mapping: `os vhost mapping create` / `POST /api/v1/vhost/mapping/`
- Issue a TLS certificate: `os ssl create-trusted` / `POST /api/v1/ssl/trusted/`
- Create a database: `os db create` / `POST /api/v1/database/:dbType/`
- Schedule a job: `os cron create` / `POST /api/v1/cron/`
- Update PHP: `os runtime php update` / `PUT /api/v1/runtime/php/:hostname/`
- Manage files: `docker cp` / `/api/v1/files/`
- Manage accounts: `os account` / `/api/v1/account/`

## Guardrails

- Prefer the CLI — run it directly when you are inside the container, or through `docker exec` when outside. Use the API only when the container shell is out of reach.
- Confirm with the user before any delete: `mktplace delete`, `services delete`, `db delete`, `db delete-user`, `account delete`, `ssl delete`, `cron delete`. Deletion destroys data and cannot be undone.
- Never write passwords or tokens into a file, a log, or a commit.
- The CLI runs inside the container as root and asks no authentication. Treat every command as a privileged action.
- OS owns the web server and the PHP configuration. Change mappings, certificates, and runtime settings through the interfaces. Hand edits to the generated files are overwritten.
- An API token expires. Log in again and repeat the call when the API answers `unauthorized`.
