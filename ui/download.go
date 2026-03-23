package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aditya-gupta-dev/oko/config"
)

func (app *App) DownloadSelectedYoutubeSong() {
	selectedSong, ok := app.widgets.GetSelectedYoutubeSong()
	if !ok {
		return
	}

	cacheDir := config.GetCacheSongsDir()
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		app.widgets.SetStatusText("Failed to create cache directory")
		return
	}

	outputTemplate := filepath.Join(cacheDir, sanitizeDownloadName(selectedSong.Title)+".%(ext)s")
	videoLink := fmt.Sprintf("https://www.youtube.com/watch?v=%s", selectedSong.VideoId)

	app.widgets.OpenDownloadDialog(selectedSong.Title)
	app.application.SetFocus(app.widgets.downloadLogs)
	app.widgets.SetStatusText(fmt.Sprintf("Downloading %q", selectedSong.Title))
	app.widgets.AppendDownloadLog(fmt.Sprintf("$ %s %s -x --audio-format mp3 -o %s", config.GetYTDlpPath(), videoLink, outputTemplate))
	app.widgets.AppendDownloadLog("")

	go app.runYoutubeDownload(config.GetYTDlpPath(), videoLink, outputTemplate, selectedSong.Title)
}

func (app *App) runYoutubeDownload(binaryPath, videoLink, outputTemplate, songTitle string) {
	cmd := exec.Command(binaryPath, videoLink, "-x", "--audio-format", "mp3", "-o", outputTemplate)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		app.finishDownloadWithError(songTitle, fmt.Errorf("failed to capture yt-dlp stdout: %w", err))
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		app.finishDownloadWithError(songTitle, fmt.Errorf("failed to capture yt-dlp stderr: %w", err))
		return
	}

	if err := cmd.Start(); err != nil {
		app.finishDownloadWithError(songTitle, fmt.Errorf("failed to start yt-dlp: %w", err))
		return
	}

	done := make(chan struct{}, 2)

	go app.streamDownloadLogs(stdout, done)
	go app.streamDownloadLogs(stderr, done)

	<-done
	<-done

	err = cmd.Wait()

	app.application.QueueUpdateDraw(func() {
		if err != nil {
			app.widgets.SetStatusText("Download failed")
			app.widgets.FinishDownloadDialog(fmt.Sprintf("Download failed for %q: %s", songTitle, err.Error()))
			return
		}

		app.widgets.SetStatusText(fmt.Sprintf("Downloaded %q", songTitle))
		app.widgets.FinishDownloadDialog(fmt.Sprintf("Download completed for %q.", songTitle))
		go app.widgets.songsList.ReloadSongs(app.application)
	})
}

func (app *App) streamDownloadLogs(reader io.Reader, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	scanner := bufio.NewScanner(reader)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		app.application.QueueUpdateDraw(func() {
			app.widgets.AppendDownloadLog(line)
		})
	}

	if err := scanner.Err(); err != nil {
		app.application.QueueUpdateDraw(func() {
			app.widgets.AppendDownloadLog(fmt.Sprintf("log stream error: %s", err.Error()))
		})
	}
}

func (app *App) finishDownloadWithError(songTitle string, err error) {
	app.application.QueueUpdateDraw(func() {
		app.widgets.SetStatusText("Download failed")
		app.widgets.FinishDownloadDialog(fmt.Sprintf("Download failed for %q: %s", songTitle, err.Error()))
	})
}

func sanitizeDownloadName(name string) string {
	sanitized := strings.TrimSpace(strings.ToLower(name))
	sanitized = strings.ReplaceAll(sanitized, "/", "-")
	sanitized = strings.ReplaceAll(sanitized, "\\", "-")
	sanitized = strings.ReplaceAll(sanitized, ":", "-")
	sanitized = strings.ReplaceAll(sanitized, "*", "")
	sanitized = strings.ReplaceAll(sanitized, "?", "")
	sanitized = strings.ReplaceAll(sanitized, "\"", "")
	sanitized = strings.ReplaceAll(sanitized, "<", "")
	sanitized = strings.ReplaceAll(sanitized, ">", "")
	sanitized = strings.ReplaceAll(sanitized, "|", "")
	sanitized = strings.ReplaceAll(sanitized, " ", "-")

	if sanitized == "" {
		return "youtube-audio"
	}

	return sanitized
}
