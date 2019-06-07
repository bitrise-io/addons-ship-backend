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

_Note: if you generate new migration while you are running dev console, the latest migration won't be executed with the command above. To do so, you can run `go run main.go up`_

## Database seeding

For having proper development data locally, you have to seed your database. There's a seeding script in the [db/seed/main.go](https://github.com/bitrise-io/addons-ship-backend/tree/master/db/seed/main.go) file. This reads the [test_data.yml](https://github.com/bitrise-io/addons-ship-backend/tree/master/db/seed/test_data.yml) file, parses it and creates the records in the development database. You can add additional data to this file and re-run the script, which will create the new ones also. In this case pay attention for the IDs of the objects, with those fields you can specify the connection between them.

_Note: To be able to generate AWS presigned URLs, you have to set the related environment variables(`AWS_BUCKET`, `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`) in your .bitrise.secrets.yml_
