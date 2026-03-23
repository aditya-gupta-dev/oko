# Oko

A terminal music player for local MP3 libraries with built-in YouTube search and download support.

![Screenshot](./screenshot.png)

## How To Use

### 1. Install dependencies

Oko expects the following tools to be available:

- Go `1.24+`
- `yt-dlp`
- `ffmpeg` or another backend supported by `yt-dlp` for audio extraction

### 2. Create the config file

On first run, Oko creates a config file at:

```text
~/.config/oko/config.json
```

You can also create it manually:

```json
{
 "folders": [
  "/absolute/path/to/your/music"
 ],
 "yt-dlp-path": "",
 "yt-api-key": ""
}
```

Configuration fields:

- `folders`: absolute paths to folders containing `.mp3` files
- `yt-dlp-path`: optional explicit path to the `yt-dlp` binary; leave empty to use `yt-dlp` from `PATH`
- `yt-api-key`: YouTube Data API key used for in-app search

### 3. Run the app

```bash
go run .
```

Or build it first:

```bash
go build -o oko .
./oko
```

### 4. Navigate

Main library controls:

- `j`: move down
- `k`: move up
- `Enter`: play selected song
- `Space`: pause or resume playback
- `a`: seek backward 5 seconds
- `d`: seek forward 5 seconds
- `/`: open YouTube search

YouTube search controls:

- Type a query and press `Enter` to search
- `Tab`: switch focus between the search input and result list
- `j` / `k`: move through YouTube results when the result list is focused
- `Enter`: download the selected YouTube result as audio
- `Esc`: close the search dialog

Download dialog:

- Shows live `yt-dlp` logs
- Stays focused until the download process finishes
- `Esc` closes the dialog only after the process is complete

## Features

- Fast local `.mp3` discovery with concurrent scanning
- Keyboard-first terminal UI built with `tview`
- Local playback with pause, resume, and seek support
- YouTube song search from inside the app
- Audio download flow powered by `yt-dlp`
- Automatic refresh of the song list after a successful download
- Cache-backed downloads stored outside your main library folders

## Project Structure

```text
.
├── api/
│   └── yt.go              # YouTube Data API search integration
├── config/
│   ├── config.go          # Config accessors and app paths
│   └── utils.go           # Config bootstrap and validation
├── song/
│   ├── player.go          # Audio playback and seek control
│   └── song.go            # MP3 metadata and file discovery
├── ui/
│   ├── download.go        # yt-dlp execution and log streaming
│   ├── songlist.go        # Library list loading and refresh
│   ├── ui.go              # App lifecycle and key bindings
│   └── widgets.go         # Shared UI widgets and dialogs
├── main.go                # Application entry point
├── go.mod
├── go.sum
└── README.md
```

## Configuration And Storage

### Config location

```text
~/.config/oko/config.json
```

### Download cache location

Downloaded audio is stored in the user cache directory:

```text
~/.cache/oko/
```

Oko reloads songs from both:

- the folders listed in `config.json`
- the Oko cache directory

This means songs downloaded from YouTube become available in the library after the download completes.

## YouTube Download Behavior

When you press `Enter` on a YouTube search result, Oko:

1. closes the search dialog
2. opens a dedicated download log dialog
3. runs `yt-dlp` against the selected video
4. extracts audio as `.mp3`
5. stores the result in the Oko cache directory
6. refreshes the in-app song list

The current implementation uses the equivalent of:

```bash
yt-dlp "https://www.youtube.com/watch?v=<video-id>" -x --audio-format mp3 -o "<cache-dir>/<sanitized-song-name>.%(ext)s"
```

## Development

### Run tests

```bash
go test ./...
```

### Format code

```bash
gofmt -w .
```

### Useful local workflow

```bash
go run .
```

## Troubleshooting

### No songs appear

- Make sure `folders` contains absolute paths
- Make sure the target directories contain `.mp3` files
- Confirm the files are readable by your user

### YouTube search fails

- Check that `yt-api-key` is set in `~/.config/oko/config.json`
- Verify the API key is valid for the YouTube Data API

### Downloads fail

- Make sure `yt-dlp` is installed and available on `PATH`
- If it is installed in a custom location, set `yt-dlp-path`
- Make sure `ffmpeg` is installed so `yt-dlp` can extract audio as MP3

### Playback fails

- Oko currently scans for and plays `.mp3` files
- If a downloaded file did not convert successfully to MP3, it will not appear in the library

## Architecture Notes

- The app is intentionally small and package-oriented
- UI state is owned inside the `ui` package
- Config setup is centralized in `config`
- Local file discovery and playback stay isolated in `song`
- External network and download integrations are separated into `api` and `ui/download.go`

## Roadmap Ideas

- Playlist support
- Better metadata display
- Background download queue
- Repeat and shuffle modes
- Search/filter for local libraries
- Cross-platform packaging

## Contributing

Issues and pull requests are welcome. If you plan to make a larger change, opening an issue first is a good way to align on direction.

## License

No license file is currently included in this repository. If you intend to open-source Oko publicly, adding a license should be one of the next steps.
