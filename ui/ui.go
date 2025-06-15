package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	application *tview.Application
	widgets     *Widgets
}

func initApp() *App {
	app := tview.NewApplication()
	widgets := InitWidgets()

	go widgets.songsList.AddSongs(app)

	return &App{
		application: app,
		widgets:     widgets,
	}
}

func StartApplication() {
	app := initApp()

	app.AttachKeyListener()

	app.application.
		SetRoot(app.widgets.rootFlex, true).
		EnableMouse(false).
		Run()
}

func (app *App) AttachKeyListener() {
	app.application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			currentIndex := app.widgets.songsList.songList.GetCurrentItem()
			itemCount := app.widgets.songsList.songList.GetItemCount()

			if currentIndex < itemCount {
				app.widgets.songsList.songList.SetCurrentItem(currentIndex + 1)
			}
			return nil
		case 'k':
			currentIndex := app.widgets.songsList.songList.GetCurrentItem()
			itemCount := app.widgets.songsList.songList.GetItemCount()

			if currentIndex < itemCount {
				app.widgets.songsList.songList.SetCurrentItem(currentIndex - 1)
			}
			return nil
		case ' ':
			if app.widgets.player.IsPlaying() {
				app.widgets.player.Pause()
				app.widgets.SetStatusText("Paused")
			} else {
				app.widgets.player.Play()
				app.widgets.SetStatusText("Playing")
			}

		}

		switch event.Key() {
		case tcell.KeyEnter:
			currentIndex := app.widgets.songsList.songList.GetCurrentItem()

			if app.widgets.player != nil {
				app.widgets.player.Cleanup()
			}

			if app.widgets.player.IsPlaying() {
				app.widgets.player.Stop()
			}

			err := app.widgets.player.LoadFile(app.widgets.songsList.songs[currentIndex].Path)

			if err != nil {
				panic(err)
			}

			app.widgets.player.Play()
			app.widgets.SetStatusText("Playing")
		}

		return event
	})
}
