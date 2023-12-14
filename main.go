package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	mainWindow = MainWindow{
		Title: "",
		MinSize: Size{
			Width:  600,
			Height: 400},
		Layout: HBox{},
	}
	searchComposite = NewSearchTabChild()
	resultsComposite = NewResultsTabChild()
	downloadsComposite = NewDownloadsTabChild()
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
var downloadsTableView *walk.TableView
var downloadsTableModel *DownloadModel

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
						Font: Font{Family: "Segoe UI", PointSize: 16},
						OnKeyPress: walk.KeyEventHandler(func(key walk.Key) {
							if key == walk.KeyReturn {
								OnSearch()
							}
						}),
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
						Text:      "Search",
						OnClicked: OnSearch,
					},
					HSpacer{},
				},
			},
			VSpacer{},
		},
	}
}

// Results tab Composite
func NewResultsTabChild() Composite {
	resultsTableModel = newResultModel()
	return Composite{
		Layout: VBox{},
		Children: []Widget{
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
									Width:  180,
									Height: 90,
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
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						MaxSize: Size{
							Width:  100,
							Height: 30,
						},
						Text: "Check all",
						OnClicked: func() {
							resultsTableModel.CheckAll()
						},
					},
					PushButton{
						MaxSize: Size{
							Width:  100,
							Height: 30,
						},
						Text: "Download selected",
						OnClicked: func() {
							var res = resultsTableModel.GetAllChecked()
							if len(res) > 0 {
								downloadsTableModel.SetDownloadRowsFromResults(res)
							}
						},
					},
					HSpacer{},
				},
			},
		},
	}
}

// Downloads tab Composite
func NewDownloadsTabChild() Composite {
	downloadsTableModel = newDownloadModel()
	return Composite{
		Layout: VBox{},
		Children: []Widget{
			TableView{
				AssignTo:         &downloadsTableView,
				AlternatingRowBG: true,
				CheckBoxes:       true,
				ColumnsOrderable: true,
				MultiSelection:   true,
				Font:             Font{Family: "Segoe UI", PointSize: 16},
				Columns: []TableViewColumn{
					{Title: ""},
					{Title: "Progress", Width: 200},
					{Title: "Title", Width: 200},
					{Title: "Channel name", Width: 200},
					{Title: "Thumbnail", Alignment: AlignCenter, Width: 160},
				},
				CustomRowHeight: 90,
				StyleCell: func(style *walk.CellStyle) {
					item := downloadsTableModel.items[style.Row()]
					switch style.Col() {
					// Thumbnail column
					case 4:
						if canvas := style.Canvas(); canvas != nil {
							bounds := style.Bounds()
							err := canvas.DrawImageStretchedPixels(item.Image,
								walk.Rectangle{
									X:      bounds.X,
									Y:      bounds.Y,
									Width:  180,
									Height: 90,
								})
							if err != nil {
								walk.MsgBox(nil, "Error", "Unable to draw walk image to canvas, "+item.ThumbnailUrl+",\nerror:\n"+err.Error(), walk.MsgBoxOK)
							}
						}
					}
				},
				Model: downloadsTableModel,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						MaxSize: Size{
							Width:  100,
							Height: 30,
						},
						Text: "Cancel all",
						// TODO: Make this
						OnClicked: func() {},
					},
					HSpacer{},
				},
			},
		},
	}
}
