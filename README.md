# Audio Gonverter

## Development 

Suggested `pre-commit` hook:

```bash
#!/usr/bin/env bash

go fmt cmd/*
go test -coverprofile cover.out -v ./...
go tool cover -func=cover.out
go tool cover -html=cover.out -o cover.html
```

Suggested `commit-msg` hook:

```bash
#!/usr/bin/env bash

if [ -z "$1" ]; then
	echo "Missing argument (commit message). Did you try to run this manually?"
	exit 1
fi

commitTitle="$(cat $1 | head -n1)"

# ignore merge requests
if echo "$commitTitle" | grep -qE "Merge branch"; then
	echo "Commit hook: ignoring branch merge"
	exit 0
fi

# check semantic versioning scheme
if ! echo "$commitTitle" | grep -qE '^(?:|feat|fix|docs|style|refactor|perf|test|chore)\(?(?:\w+|\s|\-|_)?\)?:\s\w+'; then
	echo "Your commit title did not follow semantic versioning: $commitTitle"
	echo "Please see https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commit-message-format"
	exit 1
fi
```

Run locally with:
```bash
source .env
go run cmd/main.go cmd/helpers.go cmd/routes.go
```

Or using docker-compose:
```bash
docker compose up --build
```