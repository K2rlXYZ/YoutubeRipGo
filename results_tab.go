package main

import (
	"image"
	"image/jpeg"
	"net/http"
	"strings"

	youtube "github.com/KarlMul/youtubeGo"
	"github.com/lxn/walk"
)

type Result struct {
	Index        int
	checked      bool
	Title        string
	Description  string
	ChannelTitle string
	ID           string
	ThumbnailUrl string
	Image        walk.Image
}

type ResultModel struct {
	walk.TableModelBase
	walk.SorterBase
	items []*Result
}

func newResultModel() *ResultModel {
	m := new(ResultModel)
	m.items = make([]*Result, 0)
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
		return item.Index

	case 1:
		return item.Title

	case 2:
		return item.Description

	case 3:
		return item.ChannelTitle

	case 4:
		return item.ID

	case 5:
		return item.ThumbnailUrl
	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *ResultModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *ResultModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

func (m *ResultModel) SetResultRowsFromQueryResponseData(results *youtube.QueryResponseData) {
	amount := len(results.Items)
	m.items = make([]*Result, amount)

	for i := 0; i < amount; i++ {
		m.items[i] = &Result{
			Index:        i,
			ID:           results.Items[i].ID.VideoID,
			Title:        results.Items[i].Snippet.Title,
			Description:  results.Items[i].Snippet.Description,
			ChannelTitle: results.Items[i].Snippet.ChannelTitle,
			ThumbnailUrl: results.Items[i].Snippet.Thumbnails.Medium.URL,
			Image:        nil,
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
		Index:        0,
		ID:           video.ID,
		Title:        video.Title,
		Description:  video.Description,
		ChannelTitle: video.Author,
		ThumbnailUrl: jpgUrl,
		Image:        nil,
	}

	m.PublishRowsReset()
}

func (m *ResultModel) SetResultRowsFromVideos(videos []*youtube.PlaylistEntry) {
	m.items = make([]*Result, len(videos))

	for i, video := range videos {
		var jpgUrl = strings.Replace(
			strings.Replace(
				video.Thumbnails[1].URL,
				"vi_webp",
				"vi",
				1),
			".webp",
			".jpg",
			1)

		m.items[i] = &Result{
			Index:        i,
			ID:           video.ID,
			Title:        video.Title,
			Description:  "",
			ChannelTitle: video.Author,
			ThumbnailUrl: jpgUrl,
			Image:        nil,
		}
	}

	m.PublishRowsReset()
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
