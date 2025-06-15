package song

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tcolgate/mp3"
)

type Song struct {
	Name     string
	Path     string
	Duration time.Duration
}

func CreateSongFile(path string) (Song, error) {
	file, err := os.Open(path)

	if err != nil {
		return Song{}, err
	}

	defer file.Close()

	decoder := mp3.NewDecoder(file)

	totalDuration := 0.0
	var frame mp3.Frame
	var skipped int

	for {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}

			return Song{}, err
		}
		totalDuration += frame.Duration().Seconds()
	}

	var song Song = Song{
		Path:     file.Name(),
		Name:     filepath.Base(file.Name()),
		Duration: time.Second * time.Duration(totalDuration),
	}

	return song, nil
}

func ListSongFiles(path string) ([]Song, error) {
	items, err := os.ReadDir(path)
	songs := make([]Song, 0, 20)

	if err != nil {
		return []Song{}, err
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if strings.HasSuffix(strings.ToLower(item.Name()), ".mp3") {
			song, err := CreateSongFile(filepath.Join(path, item.Name()))

			if err != nil {
				return songs, err
			}

			songs = append(songs, song)
		}
	}

	return songs, nil
}
