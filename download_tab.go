package main

import (
	"fmt"

	"github.com/lxn/walk"
)

type Download struct {
	Index        int
	Title        string
	ChannelTitle string
	ThumbnailUrl string
	Image        walk.Image
	Downloading  bool
}

type DownloadModel struct {
	walk.TableModelBase
	walk.SorterBase
	items []*Download
}

func newDownloadModel() *DownloadModel {
	m := new(DownloadModel)
	m.items = make([]*Download, 1)
	m.PublishRowsReset()
	return m
}

func (m *DownloadModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *DownloadModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 1:
		return item.Index

	case 2:
		return item.Title

	case 3:
		return item.ChannelTitle

	case 4:
		return item.ThumbnailUrl
	}

	panic("unexpected col")
}

func (m *DownloadModel) SetDownloadRowsFromResults(results []*Result) {
	fmt.Println(len(m.items))
	m.items = make([]*Download, len(results))

	for i, result := range results {
		m.items[i] = &Download{
			Index:        i,
			Title:        result.Title,
			ChannelTitle: result.ChannelTitle,
			ThumbnailUrl: result.ThumbnailUrl,
			Image:        result.Image,
			Downloading:  true,
		}
	}

	m.PublishRowsReset()
}
