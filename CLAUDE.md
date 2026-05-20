# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build (CGO is required for sqlite3)
CGO_ENABLED=1 go build -o sinister .

# Build with version embedded
CGO_ENABLED=1 go build -ldflags "-X main.VERSION=v0.2.0" -o sinister .

# Format
go fmt ./...

# Vet
go vet ./...

# Test
go test ./...

# Test a single package
go test ./manage/...

# Test a single test
go test ./manage/... -run TestFunctionName

# Release builds (linux/amd64 only if no cross-compilers available)
TARGETS="linux/amd64" bash build-release.sh v0.X.Y
```

CGO cannot be disabled — `github.com/mattn/go-sqlite3` is a hard CGO dependency.

## Architecture

`sinister` is a CLI tool that tracks YouTube subscriptions via RSS and downloads videos locally. The data flow is: **RSS feed → SQLite database → download to disk**.

### Package responsibilities

- **`main.go`** — entry point; passes version string into `argparse.Run`
- **`argparse/`** — CLI definition via `kong`. Each subcommand is a struct with a `Run(*Opts)` method. `Opts` carries global flags (`--config`, version). This is the only place where command-line flags are defined.
- **`config/`** — reads/writes the TOML config (`videoFolder`, `urls`, `quality`). Config path resolves via XDG if not overridden. `ParseConfig` is called at the top of most `manage/` functions that need the video folder or feed URLs.
- **`feed/`** — RSS XML types (`Feed`, `Entry`, `Author`, `Link`). `Entry.FilePath` is populated from the database (not from the feed itself).
- **`database/`** — all SQLite access. `openDB()` is called on every public function (no persistent connection). Schema lives in `database.go`; `migrateDB()` runs on every open and adds new columns idempotently. Queries are in `query.go`, inserts/updates in `insert.go` and `database.go`.
- **`manage/`** — one file per subcommand, each implementing its business logic by composing `config/`, `database/`, and external calls. No manage file imports another manage file; shared helpers (`makeSet`, `parsePublished`, `outputBaseName`, etc.) live in the file that needs them or in `download.go` if shared.
- **`utils/`** — XDG path helpers (largely unused now; database path is resolved directly in `database/database.go` via `xdg.DataFile`).

### Key data flows

**`update`**: fetches RSS XML over HTTP → parses into `feed.Entry` → `database.Add` inserts new entries (deduped by `slug` unique constraint).

**`download`**: queries `entries WHERE watched=0` → filters by days/count → for each entry calls `kkdai/youtube` to get video metadata → selects best format → downloads stream(s) → optionally merges with ffmpeg → marks `watched=1` and records `filepath` in DB.

**`delete`**: queries `entries WHERE filepath IS NOT NULL` → filters by creator/slug/days → `os.Remove` the file → clears `filepath` in DB → optionally marks `watched=0`.

**`sync`**: runs `Update` then `Download` then (if `--days > 0`) `Delete` in sequence.

### Database schema

Single table `entries`:
```
id INTEGER PRIMARY KEY
author TEXT
title TEXT UNIQUE
published TEXT        -- RFC3339
link TEXT
watched INTEGER       -- 0 or 1
slug TEXT UNIQUE      -- gosimple/slug of title
filepath TEXT         -- absolute path on disk, NULL if not downloaded
```

`filepath` was added via migration; all pre-migration rows have `NULL`.

### Commit message format

```
feat[codex][claude]: short description
```

Use `feat`, `fix`, `refactor`, or `chore`. Always include `[codex]` and `[claude]` tags.
