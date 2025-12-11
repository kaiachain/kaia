# How to lint your change

This document describes how to setup automatically or manually linting your change.

## Prerequisites

- `go` version should be equal to the one specified in `go.mod`

## Run linter

Run the following command (`golangci-lint` will be installed automatically if not installed):

```bash
go run build/ci.go lint -v --new-from-rev=dev
```

## Git Hook Setup

This will automatically run the linter when you commit your change. Copy and paste below pre-commit script to `.git/hooks/pre-commit` file and make the file executable (e.g., `chmod +x pre-commit`).

```bash
#!/bin/sh
go run build/ci.go lint -v --new-from-rev=dev
```
