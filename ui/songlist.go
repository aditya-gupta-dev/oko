package ui

import (
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
	folders := config.GetConfigFolders()

	for _, folder := range folders {
		songs, err := song.ListSongFilesOptimized(folder, 6) // faster with goroutines
		// songs, err := song.ListSongFiles(dir) poor performance
		if err != nil {
			panic(err)
		}

		list.songs = append(list.songs, songs...)

		if len(songs) < 1 {
			return
		}
		for _, song := range songs {
			list.songList.AddItem(song.Name, song.Duration.String(), 0, nil)
		}
		app.QueueUpdateDraw(func() {})
	}

}
