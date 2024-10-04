SHELL=/bin/bash
UI_DIR=src/presentation/ui

dev:
	air serve

build:
	templ generate -path $(UI_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/os ./os.go
	if podman ps | grep -q "os"; then podman exec os /bin/bash -c "supervisorctl restart os-api"; fi

run:
	/var/infinite/os serve