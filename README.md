[![Build Status](https://app.bitrise.io/app/dd3a27f65d8c8a5e/status.svg?token=LZLNTOHgYJwYl3TBfvmGsg&branch=master)](https://app.bitrise.io/app/dd3a27f65d8c8a5e)

# Bitrise Ship Addon Backend

## Setup

To setup the environment (database, Docker build, Go module download, etc.) run:

```
bitrise run setup
```

## Start server

To start the backend server run:

```
bitrise run up
```

## Start a development console

You can run migrations from development console, to start one run:

```
bitrise run dev-console
```

This will compile a new goose binary with the Go file-based migrations. You can run [the standard goose commands](https://github.com/pressly/goose) with the compiled `goose` binary. E.g. to migrate to the most recent migration state run:

```
./goose up
```

This is a simplified Goose CLI tool, which uses fix values and runs the original Goose CLI (for the sake of simplicity of the commands in development console). You can modify this CLI tool in the [db/main.go](https://github.com/bitrise-io/addons-ship-backend/tree/master/db/main.go) folder.
