package ui

import (
	"os"

	"github.com/aditya-gupta-dev/oko/config"
	"github.com/aditya-gupta-dev/oko/song"
	"github.com/rivo/tview"
)

const APPTITLE = " Oko - Music Player From Hell "

type SongList struct {
	songList *tview.List
	songs    []song.Song
}

func InitSongList() *SongList {
	list := tview.NewList()

	list.SetTitle(APPTITLE)
	list.SetBorder(true)

	return &SongList{
		songList: list,
	}
}

func (list *SongList) AddSongs(app *tview.Application) {
	songs, err := list.loadSongs()
	if err != nil {
		panic(err)
	}

	app.QueueUpdateDraw(func() {
		list.songs = append(list.songs, songs...)

		for _, song := range songs {
			list.songList.AddItem(song.Name, song.Duration.String(), 0, nil)
		}
	})
}

func (list *SongList) ReloadSongs(app *tview.Application) {
	songs, err := list.loadSongs()
	if err != nil {
		panic(err)
	}

	app.QueueUpdateDraw(func() {
		list.songs = songs
		list.songList.Clear()
		list.songList.SetTitle(APPTITLE)
		list.songList.SetBorder(true)

		for _, currentSong := range songs {
			list.songList.AddItem(currentSong.Name, currentSong.Duration.String(), 0, nil)
		}
	})
}

func (list *SongList) loadSongs() ([]song.Song, error) {
	folders := config.GetSongFolders()
	collectedSongs := make([]song.Song, 0, 32)

	for _, folder := range folders {
		songs, err := song.ListSongFilesOptimized(folder, 6) // faster with goroutines
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}

		collectedSongs = append(collectedSongs, songs...)
	}

	return collectedSongs, nil

}
