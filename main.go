package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"reflect"
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
var YtApiv3KEY = "AIzaSyDnbVoEtQ-nTcS-P4tIGUbSkWd2RBnpHgY"

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
							var found, res = SearchQuery(inSearchQuery.Text())
							fmt.Println(reflect.TypeOf(res))
							switch reflect.TypeOf(res) {
							case reflect.TypeOf(&youtube.Video{}):
								resultsTableModel.SetResultRowsFromVideo(res.(*youtube.Video))
							case reflect.TypeOf(&youtube.QueryResponseData{}):
								resultsTableModel.SetResultRowsFromQueryResponseData(res.(*youtube.QueryResponseData))
							case reflect.TypeOf([]*youtube.PlaylistEntry{}):
								fmt.Println(res)
								resultsTableModel.SetResultRowsFromVideos(res.([]*youtube.PlaylistEntry))
							}
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

func SearchQuery(query string) (bool, any) {
	var client = youtube.Client{
		YtApiv3Key: YtApiv3KEY,
	}

	var videoPattern, _ = regexp.Compile(`(?:watch\?v=)([\w-]*)|(?:https:\/\/youtu\.be\/)([\w-]*)(?:\?si)`)
	var matchedVideoUrl = videoPattern.FindString(query)
	if matchedVideoUrl != "" {
		var video, err = client.GetVideo(query)
		if err != nil {
			walk.MsgBox(nil, "Error", "Unable to search query, video search error:\n"+err.Error(), walk.MsgBoxOK)
			return false, nil
		}
		return true, video
	}

	var playlistPattern, _ = regexp.Compile(`(?:playlist\?list=)([\w-]*)(?:&si|$)`)
	var matchedplaylistUrl = playlistPattern.FindString(query)
	if matchedplaylistUrl != "" {
		var playlist, err = client.GetPlaylist(query)
		if err != nil {
			walk.MsgBox(nil, "Error", "Unable to search query, playlist search error:\n"+err.Error(), walk.MsgBoxOK)
			return false, nil
		}
		return true, playlist.Videos
	}

	res, err := client.SearchWithQuery(query)
	if err != nil {
		walk.MsgBox(nil, "Error", "Unable to search query, error:\n"+err.Error(), walk.MsgBoxOK)
		return false, nil
	}
	return true, res
}

func ImageFromURL(url string) image.Image {
	response, err := http.Get(url)
	if err != nil {
		walk.MsgBox(nil, "Error", "Unable to download thumbnail, "+url+",\nerror:\n"+err.Error(), walk.MsgBoxOK)
	}

	if response.StatusCode != 200 {
		walk.MsgBox(nil, "Error", "Didn't recieve 200 response code when downloading thumbnail, "+url, walk.MsgBoxOK)
	}
	defer response.Body.Close()

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		walk.MsgBox(nil, "Error", "Unable to decode thumbnail, "+url+"\nerror:\n"+err.Error(), walk.MsgBoxOK)
	}
	return img
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
				Font:             Font{Family: "Segoe UI", PointSize: 16},
				Columns: []TableViewColumn{
					{Title: ""},
					{Title: "Title", Width: 200},
					{Title: "Description", Width: 300},
					{Title: "Channel name", Width: 200},
					{Title: "ID", Width: 200},
					{Title: "Thumbnail", Alignment: AlignCenter, Width: 160},
				},
				CustomRowHeight: 90,
				StyleCell: func(style *walk.CellStyle) {
					item := resultsTableModel.items[style.Row()]
					if item.checked {
						if style.Row()%2 == 0 {
							style.BackgroundColor = walk.RGB(159, 215, 255)
						} else {
							style.BackgroundColor = walk.RGB(143, 199, 239)
						}
					}
					switch style.Col() {
					case 0:
						style.TextColor = style.BackgroundColor
					// Thumbnail column
					case 5:
						if item.Image == nil {
							img := ImageFromURL(item.ThumbnailUrl)
							var err error
							var bitmap *walk.Bitmap
							bitmap, err = walk.NewBitmapFromImageForDPI(img, resultsTableView.DPI())
							if err != nil {
								walk.MsgBox(nil, "Error", "Unable to create bitmap from image, "+item.ThumbnailUrl+",\nerror:\n"+err.Error(), walk.MsgBoxOK)
							}
							item.Image, err = walk.ImageFrom(bitmap)
							if err != nil {
								walk.MsgBox(nil, "Error", "Unable to create walk image from bitmap, "+item.ThumbnailUrl+",\nerror:\n"+err.Error(), walk.MsgBoxOK)
							}
						}
						if canvas := style.Canvas(); canvas != nil {
							bounds := style.Bounds()

							err := canvas.DrawImageStretchedPixels(item.Image,
								walk.Rectangle{
									X:      bounds.X,
									Y:      bounds.Y,
									Width:  item.Image.Size().Width / 2,
									Height: item.Image.Size().Height / 2,
								})
							if err != nil {
								walk.MsgBox(nil, "Error", "Unable to draw walk image to canvas, "+item.ThumbnailUrl+",\nerror:\n"+err.Error(), walk.MsgBoxOK)
							}
						}
					}
				},
				Model: resultsTableModel,
				OnSelectedIndexesChanged: func() {
					fmt.Printf("SelectedIndexes: %v\n", resultsTableView.SelectedIndexes())
				},
			},
			VSpacer{},
		},
	}
}
