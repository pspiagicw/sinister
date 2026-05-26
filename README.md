# `sinister`

`sinister` is a tool to track and download videos from YouTube.

- [sinister](#sinister)
    - [features](#features)
    - [installation](#installation)
    - [config](#config)
    - [usage](#usage)
        - [update](#update)
        - [status](#status)
        - [download](#download)
        - [mark](#mark)
        - [delete](#delete)
        - [sync](#sync)
        - [add](#add)
        - [clean](#clean)
        - [list](#list)
        - [export](#export)
    - [RSS Feeds](#youtube-rss-feeds)
    - [Contributing](#contributing)
    - [Disclaimer](#disclaimer)

# Features

- Tracks your subscriptions using RSS feeds and downloads whatever you are interested in.
- Remove the web interface and watch videos in your favorite video player.
- No recommendations, no ads, no distractions.
- Treat YouTube like a news feed.

# Installation

You can install `sinister` by downloading a binary from the [releases](https://github.com/pspiagicw/sinister/releases) page.

Or if you have the `Go` compiler installed:

```sh
go install github.com/pspiagicw/sinister@latest
```

> Note: `sinister` requires CGO (for SQLite). Make sure you have a C compiler (`gcc`) available.

If you use [`gox`](https://github.com/pspiagicw/gox) to manage binary packages:

```
gox install github.com/pspiagicw/sinister@latest
```

# Config

To start using `sinister`, create a config file at `~/.config/sinister/config.toml`

- Use `--config` to pass an alternate config file location.

```toml
videoFolder = "~/Videos"

urls = [
  "https://www.youtube.com/feeds/videos.xml?channel_id=UCeeFfhMcJa1kjtfZAGskOCA",
  "https://www.youtube.com/feeds/videos.xml?channel_id=UCdBK94H6oZT2Q7l0-b0xmMg"
]

quality = "hd720"
```

- `urls` — RSS feeds of the channels you want to track.
- `videoFolder` — where downloaded videos are saved.
- `quality` — minimum quality for downloads (e.g. `hd720`, `1080p`). Downloads always use the highest available quality at or above this threshold; `ffmpeg` is used to merge separate video and audio streams when needed.

> Channel URLs don't work — only RSS feed URLs work. See [YouTube RSS Feeds](#youtube-rss-feeds) for how to find them.

# Usage

### `update`

Fetches the latest videos from all configured RSS feeds and inserts new entries into the database.

```sh
sinister update
sinister update --since-days 7        # only entries from the last 7 days
sinister update --url <feed-url>      # fetch a specific feed instead of config
sinister update --limit 5            # cap at 5 entries per feed
sinister update --dry-run            # preview without writing
sinister update --retries 3          # retry failed fetches
```

![update](./gifs/update.gif)

### `status`

Shows the state of the database — total entries, watched/unwatched counts, and per-creator breakdowns.

```sh
sinister status
sinister status --creator "Channel Name"   # filter to one creator
sinister status --json                     # machine-readable output
```

![status](./gifs/status.gif)

### `download`

Downloads unwatched videos to `videoFolder`. Videos are downloaded in the highest quality available (1080p or better). If the best format has no audio, `ffmpeg` is used to merge a separate audio stream. Videos shorter than 2 minutes or longer than 1 hour are skipped automatically. YouTube Shorts are also skipped.

```sh
sinister download
sinister download --days 7       # only videos from the last 7 days
sinister download --videos 5     # download at most 5 videos
```

After downloading, each video is marked as watched and its file path is recorded in the database.

![download](./gifs/download.gif)

### `mark`

Mark videos as watched (or reset them to unwatched).

```sh
sinister mark --all-unwatched              # mark all unwatched entries as watched
sinister mark --creator "Channel Name"     # mark all entries for one creator
sinister mark --slug <slug>                # mark a specific entry by slug
sinister mark --mark-all-unwatched         # reset every entry in the DB to unwatched
sinister mark --dry-run                    # preview without writing
```

![mark](./gifs/mark.gif)

### `delete`

Delete downloaded video files from disk. The file path is cleared from the database after deletion. Use `--mark-unwatched` to reset the entry so it can be re-downloaded later.

```sh
sinister delete --days 30                   # delete videos older than 30 days
sinister delete --creator "Channel Name"    # delete all downloads for a creator
sinister delete --slug <slug>               # delete a specific video by slug
sinister delete --mark-unwatched            # reset to unwatched after deleting
sinister delete --dry-run                   # preview without removing files
```

### `sync`

Runs `update` → `download` → `delete` in sequence. Designed for use in a cron job or systemd timer to keep your library up to date automatically.

```sh
sinister sync
sinister sync --days 7      # sync last 7 days; also delete videos older than 7 days
sinister sync --videos 5    # download at most 5 new videos
sinister sync --dry-run     # dry-run the update phase
```

When `--days` is set, `sync` automatically deletes downloaded files older than that threshold after downloading.

### `add`

Adds a channel's RSS feed URL to your config file. Pass any YouTube video URL from the channel and `sinister` will look up the channel feed automatically.

```sh
sinister add "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

### `clean`

Checks every feed URL in your config and removes any that return HTTP 404. URLs that fail to connect (network errors) are kept.

```sh
sinister clean
sinister clean --dry-run    # preview which URLs would be removed
```

### `list`

Lists all channels tracked in the database along with their video counts.

```sh
sinister list
```

### `export`

Exports all unwatched video URLs to `urls.txt` and marks them as watched. Useful for feeding into external download tools.

```sh
sinister export
```

# YouTube RSS Feeds

There are multiple ways to get the RSS feed URL for a YouTube channel.

The most reliable method is to view the channel page source and search for `rss`.

You can also use `sinister add` with any video URL from the channel — it will find and add the feed automatically.

Other resources:

- [Feeder's guide to YouTube RSS](https://feeder.co/knowledge-base/rss-feed-creation/youtube-rss/)

# Contributing

If you want to contribute, open an issue or a pull request on [GitHub](https://github.com/pspiagicw/sinister).

# Disclaimer

Downloading videos from YouTube may be against their terms of service. Use at your own risk.
