# TODO

## Subcommands

- [ ] `mark`: complete implementation and add non-interactive filters.
- [ ] `list`: list videos from the local database with filters.
- [ ] `sync` (or `auto`): run update + export in one command for cron/systemd.

## Existing command improvements

- [ ] `status`:
  - [ ] `--json` output.
  - [ ] `--creator <name>` creator-specific stats.
  - [ ] `--top-creators <n>` ranking output.
  - [ ] `--since <days>` recent-only statistics.

- [ ] `update`:
  - [ ] `--url <feed-url>` (repeatable) to override configured feeds.
  - [ ] `--limit <n>` max entries processed per feed.
  - [ ] `--since <days>` ignore entries older than N days.
  - [ ] `--dry-run` to report additions without DB writes.
  - [ ] `--json` summary output (inserted/skipped/error counts).

- [ ] `export`:
  - [ ] `--output <file>` destination file path.
  - [ ] `--format txt|json|aria2`.
  - [ ] `--days <n>` export only recent unwatched videos.
  - [ ] `--creator <name>` export a single creator.
  - [ ] `--no-mark` export without changing watched state.
  - [ ] `--mark` explicit mark-after-export mode.
