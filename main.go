package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

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

var inSearchQuery *walk.TextEdit
var mainWindow MainWindow
var searchComposite Composite
var resultsComposite Composite
var downloadsComposite Composite
var tabWidget *walk.TabWidget

var tabButtonWidth int = 60

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
							// TODO: add search functionality
							// TODO: add search results to resultsComposite
							tabWidget.SetCurrentIndex(resultsTabEnum)
						},
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
	return Composite{
		Layout: VBox{},
		Children: []Widget{
			VSpacer{},

			PushButton{
				MaxSize: Size{
					Width:  100,
					Height: 30,
				},
				Text: "Search",
				OnClicked: func() {

				},
			},
			VSpacer{},
		},
	}
}
