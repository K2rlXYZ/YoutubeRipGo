package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var inSearch *walk.TextEdit

var tabButtonWidth int = 60

// Tab Composite
func TabComposite() Composite {
	return Composite{
		Layout:    HBox{},
		Alignment: AlignHNearVNear,
		Children: []Widget{
			Composite{
				Layout:    HBox{},
				Alignment: AlignHNearVNear,
				Children: []Widget{
					PushButton{
						Alignment: AlignHNearVNear,
						Text:      "Search",
						MaxSize: Size{
							Width: tabButtonWidth,
						},
					},
					PushButton{
						Alignment: AlignHNearVNear,
						Text:      "Results",
						MaxSize: Size{
							Width: tabButtonWidth,
						},
					},
					PushButton{
						Alignment: AlignHNearVNear,
						Text:      "Downloads",
						MaxSize: Size{
							Width: tabButtonWidth,
						},
					},
					/*
						PushButton{
							Alignment: AlignHNearVNear,
							Text: "Browse",
							MaxSize: Size{
								Width: tabButtonWidth,
							},
						},
					*/
				},
			},
			HSpacer{},
			HSpacer{},
		},
	}
}

func SearchWindow() Composite {
	return Composite{
		Layout: VBox{},
		Children: []Widget{
			TabComposite(),
			VSpacer{},
			//Search query Composite
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
						AssignTo: &inSearch,
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
							walk.MsgBox(nil, "Button Clicked", "You clicked the button.", walk.MsgBoxOK)
						},
					},
					HSpacer{},
				},
			},
			VSpacer{},
		},
	}
}

func main() {
	var mainWindow = MainWindow{
		Title: "SCREAMO",
		MinSize: Size{
			Width:  600,
			Height: 400},
		Layout: HBox{},
	}
	mainWindow.Children = []Widget{SearchWindow()}
	mainWindow.Run()
}
