package ui

import (
	"fmt"

	"github.com/aditya-gupta-dev/oko/song"
	"github.com/rivo/tview"
)

type Widgets struct {
	rootFlex   *tview.Flex
	songsList  *SongList
	statusText *tview.TextView
	bottomFlex *tview.Flex
	player     *song.Player
}

func InitWidgets() *Widgets {
	songList := InitSongList()

	statusText := tview.NewTextView()

	statusText.
		SetText("Select a song").
		SetBorder(true)

	bottomFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusText, 0, 3, true)

	rootFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(songList.songList, 0, 8, true).
		AddItem(bottomFlex, 0, 2, false)

	return &Widgets{
		rootFlex:   rootFlex,
		songsList:  songList,
		statusText: statusText,
		bottomFlex: bottomFlex,
		player:     song.NewPlayer(),
	}
}

func (widget *Widgets) SetStatusText(text string) {
	if widget.player.IsPlaying() {
		widget.statusText.SetText(fmt.Sprintf("[ %s ]\n%s", text, "Currently Paused"))
	} else {
		widget.statusText.SetText(fmt.Sprintf("[ %s ]\n%s", text, "Currently Playing"))
	}
}
