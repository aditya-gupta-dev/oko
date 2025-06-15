package ui

import (
	"fmt"

	"github.com/aditya-gupta-dev/oko/song"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Widgets struct {
	rootFlex     *tview.Flex
	songsList    *SongList
	statusText   *tview.TextView
	progressText *tview.TextView
	bottomFlex   *tview.Flex
	player       *song.Player
}

func InitWidgets() *Widgets {
	songList := InitSongList()

	statusText := tview.NewTextView()

	progressText := tview.NewTextView()

	progressText.SetTextColor(tcell.Color110)
	progressText.SetText("")

	statusText.
		SetText("Select a song")

	bottomFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(statusText, 0, 3, false).
		AddItem(progressText, 0, 7, false)

	bottomFlex.SetBorder(true)

	rootFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(songList.songList, 0, 8, true).
		AddItem(bottomFlex, 0, 2, false)

	return &Widgets{
		rootFlex:     rootFlex,
		songsList:    songList,
		statusText:   statusText,
		bottomFlex:   bottomFlex,
		player:       song.NewPlayer(),
		progressText: progressText,
	}
}

func (widget *Widgets) SetStatusText(text string) {
	if widget.player.IsPlaying() {
		widget.statusText.SetText(fmt.Sprintf("[ %s ]", text))
	} else {
		widget.statusText.SetText(fmt.Sprintf("[ %s ]", text))
	}
}

// TODO: Setup progress bar and actual updates
func (widget *Widgets) SetProgress(text string) {
	if widget.player.IsPlaying() {

	}
}
