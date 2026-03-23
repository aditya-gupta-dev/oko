package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/aditya-gupta-dev/oko/api"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	application *tview.Application
	widgets     *Widgets
}

const youtubeSearchResultLimit int64 = 5

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

	app.widgets.searchInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			app.SearchYoutube(app.widgets.searchInput.GetText())
		case tcell.KeyEscape:
			app.CloseSearchDialog()
		}
	})

	app.application.
		SetRoot(app.widgets.pages, true).
		EnableMouse(false).
		Run()

}

func (app *App) AttachKeyListener() {
	app.application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if app.widgets.IsSearchDialogOpen() {
			switch event.Key() {
			case tcell.KeyEscape:
				app.CloseSearchDialog()
				return nil
			case tcell.KeyTab:
				app.ToggleSearchDialogFocus()
				return nil
			case tcell.KeyEnter:
				if app.application.GetFocus() == app.widgets.searchResults {
					return nil
				}
			}

			switch event.Rune() {
			case 'j':
				if app.application.GetFocus() == app.widgets.searchResults {
					app.MoveSearchResults(1)
					return nil
				}
			case 'k':
				if app.application.GetFocus() == app.widgets.searchResults {
					app.MoveSearchResults(-1)
					return nil
				}
			}

			return event
		}

		switch event.Rune() {
		case '/':
			app.OpenSearchDialog()
			return nil
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
		case 'a':
			if app.widgets.player.IsPlaying() {
				app.widgets.player.Seek(-5)
			}
			return nil
		case 'd':
			if app.widgets.player.IsPlaying() {
				app.widgets.player.Seek(5)
			}
		case ' ':
			if app.widgets.player.IsPlaying() {
				app.widgets.player.Pause()
				app.widgets.SetStatusText("Paused")
			} else {
				app.widgets.player.Play()
				app.widgets.SetStatusText("Playing")
				go app.UpdateProgressBar()
			}
		}

		switch event.Key() {
		case tcell.KeyEnter:
			currentIndex := app.widgets.songsList.songList.GetCurrentItem()

			if app.widgets.player != nil {
				app.widgets.player.Cleanup()
				app.widgets.SetProgress("")
			}

			if app.widgets.player.IsPlaying() {
				app.widgets.player.Stop()
				app.widgets.SetProgress("")
			}

			err := app.widgets.player.LoadFile(app.widgets.songsList.songs[currentIndex].Path)

			if err != nil {
				panic(err)
			}

			app.widgets.player.Play()
			app.widgets.SetStatusText("Playing")

			go app.UpdateProgressBar()
		}

		return event
	})
}

func (app *App) OpenSearchDialog() {
	app.widgets.OpenSearchDialog()
	app.application.SetFocus(app.widgets.searchInput)
}

func (app *App) CloseSearchDialog() {
	app.widgets.CloseSearchDialog()
	app.application.SetFocus(app.widgets.songsList.songList)
}

func (app *App) ToggleSearchDialogFocus() {
	if app.application.GetFocus() == app.widgets.searchInput {
		app.application.SetFocus(app.widgets.searchResults)
		return
	}

	app.application.SetFocus(app.widgets.searchInput)
}

func (app *App) MoveSearchResults(offset int) {
	itemCount := app.widgets.searchResults.GetItemCount()
	if itemCount == 0 {
		return
	}

	currentIndex := app.widgets.searchResults.GetCurrentItem()
	nextIndex := currentIndex + offset

	if nextIndex < 0 {
		nextIndex = 0
	}

	if nextIndex >= itemCount {
		nextIndex = itemCount - 1
	}

	app.widgets.searchResults.SetCurrentItem(nextIndex)
}

func (app *App) SearchYoutube(query string) {
	trimmedQuery := strings.TrimSpace(query)

	if trimmedQuery == "" {
		app.widgets.ResetYoutubeResults()
		return
	}

	app.widgets.searchResults.Clear()
	app.widgets.searchResults.AddItem("Searching Youtube...", "", 0, nil)
	app.widgets.SetStatusText(fmt.Sprintf("Searching Youtube for %q", trimmedQuery))

	go func() {
		results, err := api.SearchSongYoutube(trimmedQuery, youtubeSearchResultLimit)

		app.application.QueueUpdateDraw(func() {
			if err != nil {
				app.widgets.searchResults.Clear()
				app.widgets.searchResults.AddItem(fmt.Sprintf("Search failed: %s", err.Error()), "", 0, nil)
				app.widgets.SetStatusText("Youtube search failed")
				return
			}

			hasResults := app.widgets.SetYoutubeResults(results, int(youtubeSearchResultLimit))
			if hasResults {
				app.application.SetFocus(app.widgets.searchResults)
			}
			app.widgets.SetStatusText(fmt.Sprintf("Youtube search completed for %q", trimmedQuery))
		})
	}()
}

func (app *App) UpdateProgressBar() {
	for {
		if !app.widgets.player.IsPlaying() {
			return
		}

		songs := app.widgets.songsList.songs
		currentIndex := app.widgets.songsList.songList.GetCurrentItem()

		if !(currentIndex < len(songs)) {
			return
		}

		player := app.widgets.player
		currentSamples, totalSamples := player.Position()
		currentDuration, totalDuration := player.PositionDuration()

		durationStatement := fmt.Sprintf("Played [ %s ] out of [ %s ]", currentDuration, totalDuration)
		sampleStatement := fmt.Sprintf("Hearing Sample No. [ %d ] out of [ %d ] Samples", currentSamples, totalSamples)

		app.widgets.SetProgress(fmt.Sprintf("%s\n%s", durationStatement, sampleStatement))

		time.Sleep(1 * time.Second)
		app.application.QueueUpdateDraw(func() {

		})
	}
}
