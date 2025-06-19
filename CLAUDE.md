# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is Frank Chiarulli Jr.'s personal GitHub profile repository. It contains a Go-based README generator that automatically updates the profile README.md file every hour via GitHub Actions.

## Architecture

The repository consists of two main components:

1. **Root README.md** - The generated profile README that appears on GitHub
2. **generate/** directory - Contains the Go application that generates the README

The Go generator (`generate/main.go`) uses an embedded template (`README.md.tmpl`) to render the profile content and outputs it to the root README.md file.

## Commands

### Build and run the generator
```bash
cd generate
go run main.go ../README.md
```

### Lint the Go code
```bash
cd generate
golangci-lint run
```

### Test the generator locally
```bash
cd generate
go run main.go test-output.md
```

## CI/CD

The repository uses GitHub Actions (`.github/workflows/ci.yml`) that:
- Runs on every push, pull request, and hourly via cron
- Lints the Go code using golangci-lint
- Generates and commits updated README.md automatically on the main branch
- Uses Go 1.19 as specified in the workflow environment

## File Structure

- `README.md` - Auto-generated profile content (do not edit manually)
- `generate/main.go` - README generator application
- `generate/README.md.tmpl` - Template for the README content
- `generate/go.mod` - Go module definition
- `.github/workflows/ci.yml` - CI/CD pipeline configuration