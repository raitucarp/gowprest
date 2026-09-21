---
title: "Installation"
description: "How to add gowprest to your Go project."
---

## Requirements

- **Go 1.22** or higher.
- A WordPress instance with the REST API enabled (WordPress 4.7+ comes with REST API built-in).

## Installing via `go get`

Add `gowprest` to your project dependencies:

```bash
go get -u github.com/raitucarp/gowprest
```

## Verifying in `go.mod`

Your `go.mod` should include:

```go
module your-project

go 1.22

require (
    github.com/raitucarp/gowprest v1.0.0
)
```

Run `go mod tidy` to download and verify dependencies:

```bash
go mod tidy
```
