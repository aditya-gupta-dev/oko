package ui

import (
	"fmt"

	"github.com/aditya-gupta-dev/oko/api"
	"github.com/aditya-gupta-dev/oko/song"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Widgets struct {
	rootFlex            *tview.Flex
	pages               *tview.Pages
	songsList           *SongList
	statusText          *tview.TextView
	progressText        *tview.TextView
	bottomFlex          *tview.Flex
	player              *song.Player
	searchInput         *tview.InputField
	searchResults       *tview.List
	searchDialog        tview.Primitive
	searchDialogVisible bool
	youtubeSongs        []api.YoutubeSong
	downloadLogs        *tview.TextView
	downloadDialog      tview.Primitive
	downloadDialogOpen  bool
	downloadInProgress  bool
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

	searchInput := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(40)

	searchInput.SetLabelColor(tcell.Color110)
	searchInput.SetFieldTextColor(tcell.ColorWhite)
	searchInput.SetFieldBackgroundColor(tcell.Color235)
	searchInput.SetBackgroundColor(tcell.Color235)

	searchResults := tview.NewList()

	searchResults.SetBorder(true)
	searchResults.SetTitle(" Youtube Results ")
	searchResults.SetMainTextColor(tcell.ColorWhite)
	searchResults.SetSecondaryTextColor(tcell.Color110)
	searchResults.SetSelectedTextColor(tcell.ColorWhite)
	searchResults.SetSelectedBackgroundColor(tcell.Color237)
	searchResults.SetHighlightFullLine(true)

	searchDialogContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(searchInput, 3, 0, true).
		AddItem(searchResults, 0, 1, false)

	searchDialogFrame := tview.NewFrame(searchDialogContent).
		SetBorders(0, 0, 0, 0, 0, 0)

	searchDialogFrame.SetBorder(true)
	searchDialogFrame.SetTitle(" Search Youtube ")
	searchDialogFrame.AddText("Press Enter to search, Esc to close.", true, tview.AlignCenter, tcell.Color110)

	searchDialog := tview.NewGrid().
		SetRows(0, 11, 0).
		SetColumns(0, 70, 0).
		AddItem(searchDialogFrame, 1, 1, 1, 1, 0, 0, true)

	downloadLogs := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetWrap(true)

	downloadLogs.SetBorder(true)
	downloadLogs.SetTitle(" Youtube Download ")
	downloadLogs.SetTextColor(tcell.ColorWhite)

	downloadDialogFrame := tview.NewFrame(downloadLogs).
		SetBorders(0, 0, 0, 0, 0, 0)

	downloadDialogFrame.SetBorder(true)
	downloadDialogFrame.SetTitle(" Download Audio ")
	downloadDialogFrame.AddText("yt-dlp logs will appear here. Esc closes after the process finishes.", true, tview.AlignCenter, tcell.Color110)

	downloadDialog := tview.NewGrid().
		SetRows(0, 18, 0).
		SetColumns(0, 90, 0).
		AddItem(downloadDialogFrame, 1, 1, 1, 1, 0, 0, true)

	pages := tview.NewPages().
		AddPage("main", rootFlex, true, true).
		AddPage("search", searchDialog, true, false).
		AddPage("download", downloadDialog, true, false)

	return &Widgets{
		rootFlex:       rootFlex,
		pages:          pages,
		songsList:      songList,
		statusText:     statusText,
		bottomFlex:     bottomFlex,
		player:         song.NewPlayer(),
		progressText:   progressText,
		searchInput:    searchInput,
		searchResults:  searchResults,
		searchDialog:   searchDialog,
		downloadLogs:   downloadLogs,
		downloadDialog: downloadDialog,
	}
}

func (widget *Widgets) SetStatusText(text string) {
	if widget.player.IsPlaying() {
		widget.statusText.SetText(fmt.Sprintf("[ %s ]", text))
	} else {
		widget.statusText.SetText(fmt.Sprintf("[ %s ]", text))
	}
}

func (widget *Widgets) SetProgress(text string) {
	if widget.player.IsPlaying() {
		widget.progressText.SetText(text)
	}
}

func (widget *Widgets) OpenSearchDialog() {
	widget.searchDialogVisible = true
	widget.searchInput.SetText("")
	widget.ResetYoutubeResults()
	widget.pages.ShowPage("search")
}

func (widget *Widgets) CloseSearchDialog() {
	widget.searchDialogVisible = false
	widget.pages.HidePage("search")
}

func (widget *Widgets) IsSearchDialogOpen() bool {
	return widget.searchDialogVisible
}

func (widget *Widgets) ResetYoutubeResults() {
	widget.youtubeSongs = nil
	widget.searchResults.Clear()
	widget.searchResults.AddItem(" Search for a song", "  Press Enter to load the first 5 matches.", 0, nil)
}

func (widget *Widgets) SetYoutubeResults(songs []api.YoutubeSong, maxResults int) bool {
	widget.searchResults.Clear()
	widget.youtubeSongs = make([]api.YoutubeSong, 0, maxResults)

	itemCount := 0

	for _, result := range songs {
		if result.Title == "" || result.VideoId == "" {
			continue
		}

		widget.searchResults.AddItem(
			fmt.Sprintf(" %d. %s", itemCount+1, result.Title),
			fmt.Sprintf("    video id: %s", result.VideoId),
			0,
			nil,
		)
		widget.youtubeSongs = append(widget.youtubeSongs, result)
		itemCount++

		if itemCount >= maxResults {
			break
		}
	}

	if itemCount == 0 {
		widget.searchResults.AddItem(" No songs found", "  Try a different search phrase.", 0, nil)
		return false
	}

	widget.searchResults.SetCurrentItem(0)

	return true
}

func (widget *Widgets) GetSelectedYoutubeSong() (api.YoutubeSong, bool) {
	currentIndex := widget.searchResults.GetCurrentItem()

	if currentIndex < 0 || currentIndex >= len(widget.youtubeSongs) {
		return api.YoutubeSong{}, false
	}

	return widget.youtubeSongs[currentIndex], true
}

func (widget *Widgets) OpenDownloadDialog(title string) {
	widget.CloseSearchDialog()
	widget.downloadDialogOpen = true
	widget.downloadInProgress = true
	widget.downloadLogs.SetTitle(fmt.Sprintf(" Youtube Download - %s ", title))
	widget.downloadLogs.SetText("")
	widget.pages.ShowPage("download")
}

func (widget *Widgets) AppendDownloadLog(text string) {
	fmt.Fprintln(widget.downloadLogs, text)
	widget.downloadLogs.ScrollToEnd()
}

func (widget *Widgets) FinishDownloadDialog(message string) {
	widget.downloadInProgress = false
	widget.AppendDownloadLog("")
	widget.AppendDownloadLog(message)
	widget.AppendDownloadLog("Press Esc to close this dialog.")
}

func (widget *Widgets) CloseDownloadDialog() {
	widget.downloadDialogOpen = false
	widget.downloadInProgress = false
	widget.pages.HidePage("download")
}

func (widget *Widgets) IsDownloadDialogOpen() bool {
	return widget.downloadDialogOpen
}

func (widget *Widgets) CanCloseDownloadDialog() bool {
	return widget.downloadDialogOpen && !widget.downloadInProgress
}
