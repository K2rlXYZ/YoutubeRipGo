package main

import (
	"fmt"
	"regexp"

	youtube "github.com/KarlMul/youtubeGo"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type MyClient struct {
	*youtube.Client
}

func main() {
	mainWindow = MainWindow{
		Title: "SCREAMO",
		MinSize: Size{
			Width:  600,
			Height: 400},
		Layout: HBox{},
	}
	searchComposite = NewSearchTabChild()
	resultsComposite = NewResultsTabChild()
	// TODO: make DownloadsComposite function
	// downloadsWindow = DownloadsComposite()
	mainWindow.Children = []Widget{
		TabWidget{
			AssignTo: &tabWidget,
			Pages: []TabPage{
				{
					Name:     "Search",
					Title:    "Search",
					Layout:   HBox{},
					Children: []Widget{searchComposite},
				},
				{
					Name:     "Results",
					Title:    "Results",
					Layout:   HBox{},
					Children: []Widget{resultsComposite},
				},
				{
					Name:     "Downloads",
					Title:    "Downloads",
					Layout:   HBox{},
					Children: []Widget{downloadsComposite},
				},
			},
		},
	}
	mainWindow.Run()
}

var mainWindow MainWindow
var tabWidget *walk.TabWidget

var searchComposite Composite
var inSearchQuery *walk.TextEdit

var resultsComposite Composite
var resultsTableView *walk.TableView
var resultsTableModel *ResultModel

var downloadsComposite Composite

// Enums have to be in the same order as tabWidget pages for tab switching to work properly
const (
	searchTabEnum    = iota
	resultsTabEnum   = iota
	downloadsTabEnum = iota
)

// Search tab Composite
func NewSearchTabChild() Composite {
	return Composite{
		Layout: VBox{},
		Children: []Widget{
			VSpacer{},
			Composite{
				Layout: HBox{},
				MaxSize: Size{
					Height: 60,
				},
				Children: []Widget{
					HSpacer{},
					TextEdit{
						ColumnSpan: 1,
						MaxSize: Size{
							Width: 1000,
						},
						Font:     Font{Family: "Segoe UI", PointSize: 16},
						AssignTo: &inSearchQuery,
					},
					HSpacer{},
				},
			},
			//Search button Composite
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						MaxSize: Size{
							Width:  100,
							Height: 30,
						},
						Text: "Search",
						OnClicked: func() {
							var found = SearchQuery(inSearchQuery.Text())
							// TODO: add search results to resultsComposite
							if found {
								tabWidget.SetCurrentIndex(resultsTabEnum)
							}
						},
					},
					HSpacer{},
				},
			},
			VSpacer{},
		},
	}
}

func SearchQuery(query string) bool {
	var client = new(youtube.Client)

	var videoPattern, _ = regexp.Compile(`(?:watch\?v=)([\w-]*)|(?:https:\/\/youtu\.be\/)([\w-]*)(?:\?si)`)
	var matchedVideoUrl = videoPattern.FindString(query)
	if matchedVideoUrl != "" {
		var video, err = client.GetVideo(query)
		if err != nil {
			walk.MsgBox(nil, "Error", "Unable to search query, video search error:\n"+err.Error(), walk.MsgBoxOK)
			return false
		}
		// TODO: add video to return
		resultsTableModel.SetResultRowsFromVideo(video)
		fmt.Println(video.Description)
		return true
	}

	var playlistPattern, _ = regexp.Compile(`(?:playlist\?list=)([\w-]*)(?:&si|$)`)
	var matchedplaylistUrl = playlistPattern.FindString(query)
	if matchedplaylistUrl != "" {
		var playlist, err = client.GetPlaylist(query)
		if err != nil {
			walk.MsgBox(nil, "Error", "Unable to search query, playlist search error:\n"+err.Error(), walk.MsgBoxOK)
			return false
		}
		// TODO: add playlist to return
		fmt.Println(playlist.Description)
		return true
	}

	// TODO implement this
	// client.SearchWithQuery(query)

	walk.MsgBox(nil, "Error", "Unable to search query", walk.MsgBoxOK)
	return false
}

// Results tab Composite
func NewResultsTabChild() Composite {
	resultsTableModel = newResultModel()
	return Composite{
		Layout: VBox{},
		Children: []Widget{
			VSpacer{},

			TableView{
				AssignTo:         &resultsTableView,
				AlternatingRowBG: true,
				CheckBoxes:       true,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Columns: []TableViewColumn{
					{Title: "ID"},
					{Title: "Title"},
					{Title: "Description"},
					{Title: "Channel name"},
					{Title: "Thumbnail"},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := resultsTableModel.items[style.Row()]
					if canvas := style.Canvas(); canvas != nil {
						fnt := resultsTableView.Font()
						bnds := style.Bounds()
						switch style.Col() {
						case 0:
							canvas.DrawTextPixels(item.ID, fnt, 0, bnds, walk.TextLeft)
						case 1:
							canvas.DrawTextPixels(item.Title, fnt, 0, bnds, walk.TextLeft)
						case 2:
							canvas.DrawTextPixels(item.Description, fnt, 0, bnds, walk.TextLeft)
						case 3:
							canvas.DrawTextPixels(item.ChannelTitle, fnt, 0, bnds, walk.TextLeft)
						case 4:
							canvas.DrawTextPixels(item.ThumbnailUrl, fnt, 0, bnds, walk.TextLeft)
						}
					}
				},
				Model: resultsTableModel,
			},
			VSpacer{},
		},
	}
}
