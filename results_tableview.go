package main

import (
	"strings"

	youtube "github.com/KarlMul/youtubeGo"
	"github.com/lxn/walk"
)

type Result struct {
	ID           string
	Title        string
	Description  string
	ChannelTitle string
	ThumbnailUrl string
}

type ResultModel struct {
	walk.TableModelBase
	walk.SorterBase
	items []*Result
}

func newResultModel() *ResultModel {
	m := new(ResultModel)
	m.items = make([]*Result, 1)
	m.items[0] = &Result{
		ID:           "a",
		Title:        "a",
		Description:  "a",
		ChannelTitle: "a",
		ThumbnailUrl: "a",
	}

	m.PublishRowsReset()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *ResultModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *ResultModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.ID

	case 1:
		return item.Title

	case 2:
		return item.Description

	case 3:
		return item.ChannelTitle

	case 4:
		return item.ThumbnailUrl
	}

	panic("unexpected col")
}

/*// Called by the TableView to retrieve if a given row is checked.
func (m *ResultModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *ResultModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}*/

func (m *ResultModel) SetResultRowsFromQueryResponseData(results *youtube.QueryResponseData) {
	amount := len(results.Items)
	m.items = make([]*Result, amount)

	for i := 0; i < amount; i++ {
		m.items[i] = &Result{
			ID:           results.Items[i].ID.VideoID,
			Title:        results.Items[i].Snippet.Title,
			Description:  results.Items[i].Snippet.Description,
			ChannelTitle: results.Items[i].Snippet.ChannelTitle,
			ThumbnailUrl: results.Items[i].Snippet.Thumbnails.Medium.URL,
		}
	}

	m.PublishRowsReset()
}

func (m *ResultModel) SetResultRowsFromVideo(video *youtube.Video) {
	m.items = make([]*Result, 1)

	var jpgUrl = strings.Replace(
		strings.Replace(
			video.Thumbnails[1].URL,
			"vi_webp",
			"vi",
			1),
		".webp",
		".jpg",
		1)

	m.items[0] = &Result{
		ID:           video.ID,
		Title:        video.Title,
		Description:  video.Description,
		ChannelTitle: video.ChannelHandle,
		ThumbnailUrl: jpgUrl,
	}

	m.PublishRowsReset()
}
