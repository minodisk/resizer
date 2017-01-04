# Contributing

Do everything in Docker.

## Ready for development

1. Create a directory `./.secret`.
2. Create the service account JSON of Google Cloud at `./.secret/gcloud.json`.

## Test

Call `./bin/test` in your terminal.

## Run

Call `./bin/run` in your terminal.

## Before writing some code

Call `./bin/glide install` to create vendor directory in your host machine at the first time before writing some code.

## Glide

Call `./bin/glide` with some options.

- `./bin/glide update`: update packages
- `./bin/glide get path/of/package`: get some packages

## Release

Merge into release branch.
