package song

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

func ListSongFilesOptimized(path string, workers int) ([]Song, error) {

	items, err := os.ReadDir(path)

	if err != nil {
		return []Song{}, err
	}

	if len(items) < 1 {
		return []Song{}, nil
	}

	songs := make([]Song, 0, len(items))
	var chunkSize int = len(items) / workers
	var remainder int = len(items) % workers

	var wg sync.WaitGroup
	var mutex sync.Mutex

	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		chunk := items[start:end]

		if len(chunk) == 0 {
			continue
		}

		wg.Add(1)
		go func(workerId int, chunk []os.DirEntry) {
			defer wg.Done()
			var tempSongs []Song = make([]Song, 0, len(chunk))

			for _, entry := range chunk {
				if !(strings.HasSuffix(strings.ToLower(entry.Name()), ".mp3")) {
					continue
				}
				song, err := CreateSongFile(filepath.Join(path, entry.Name()))
				if err != nil {
					fmt.Println("Panic from worker : ", workerId, entry.Name())
					panic(err)
				}
				tempSongs = append(tempSongs, song)
			}
			mutex.Lock()
			songs = append(songs, tempSongs...)
			mutex.Unlock()
		}(i, chunk)
	}

	if remainder > 0 {
		chunk := items[workers*chunkSize:]
		wg.Add(1)
		go func(workerId int, chunk []os.DirEntry) {
			defer wg.Done()
			var tempSongs []Song = make([]Song, 0, len(chunk))

			for _, entry := range chunk {

				if !(strings.HasSuffix(strings.ToLower(entry.Name()), ".mp3")) {
					continue
				}
				song, err := CreateSongFile(filepath.Join(path, entry.Name()))
				if err != nil {
					fmt.Println("Panic from worker : ", workerId)
					panic(err)
				}
				tempSongs = append(tempSongs, song)
			}
			mutex.Lock()
			songs = append(songs, tempSongs...)
			mutex.Unlock()
		}(10025, chunk) // 10025 -> remainder worker id
	}

	wg.Wait()

	return songs, nil
}

func ListSongFiles(path string) ([]Song, error) {
	items, err := os.ReadDir(path)
	songs := make([]Song, 0, 20)

	if err != nil {
		return []Song{}, err
	}

	if len(items) < 1 {
		return []Song{}, nil
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
