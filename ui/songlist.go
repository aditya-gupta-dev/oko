package ui

import (
	"github.com/aditya-gupta-dev/oko/song"
	"github.com/rivo/tview"
)

const APPTITLE = " Oko - Music Player From Hell "
const DIR = "C:/Users/hyper/progs/ytt/output"

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
	songs, err := song.ListSongFiles(DIR)
	if err != nil {
		panic(err)
	}

	list.songs = songs

	if len(songs) < 1 {
		list.songList.AddItem("No Songs in the directory.", "", 0, nil)
	}
	for index, song := range songs {
		list.songList.AddItem(song.Name, song.Duration.String(), 0, nil)
		if index == 1 {
			app.QueueUpdateDraw(func() {})
		}
	}
	app.QueueUpdateDraw(func() {})
}
