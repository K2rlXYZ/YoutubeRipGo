package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	youtube "github.com/KarlMul/youtubeGo"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type Download struct {
	Index        int
	ID           string
	Title        string
	ChannelTitle string
	ThumbnailUrl string
	Image        walk.Image
	Downloading  bool
	cancelFunc   context.CancelFunc
}

type DownloadModel struct {
	walk.TableModelBase
	walk.SorterBase
	items []*Download
}

func newDownloadModel() *DownloadModel {
	m := new(DownloadModel)
	m.items = make([]*Download, 0)
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
	case 0:
		return Composite{
			Layout: HBox{},
			Children: []Widget{
				PushButton{
					MaxSize: Size{
						Width:  50,
						Height: 30,
					},
					Text: "Cancel",
					OnClicked: func() {
						item.cancelFunc()
					},
				},
			},
		}
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

type readerCtx struct {
	ctx context.Context
	r   io.Reader
}

func (r *readerCtx) Read(p []byte) (n int, err error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

// NewReader gets a context-aware io.Reader.
func NewReader(ctx context.Context, r io.Reader) io.Reader {
	return &readerCtx{ctx: ctx, r: r}
}

func (m *DownloadModel) SetAndStartDownloadsFromResults(results []*Result) {
	fmt.Println(len(m.items))
	m.items = make([]*Download, len(results))

	for i, result := range results {
		m.items[i] = &Download{
			Index:        i,
			Title:        result.Title,
			ChannelTitle: result.ChannelTitle,
			ThumbnailUrl: result.ThumbnailUrl,
			ID:           result.ID,
			Image:        result.Image,
		}
		//Async download
		go func(result *Result, down *Download) {

			//Make a cancel function, this can be later called to cancel the downloading, and add it to the download
			ctx, cancel := context.WithCancel(context.Background())
			down.cancelFunc = cancel

			execPath, _ := os.Executable()
			/*downloader := downloader.Downloader{
				OutputDir: execPath + "\\..\\videos\\",
			}
			video, _ := downloader.GetVideo(result.ID)
			video.Formats.WithAudioChannels().FindByQuality()
			downloader.Download(ctx, video, video.Forma)*/

			cli := youtube.Client{}
			vid, _ := cli.GetVideo(result.ID)
			reader, _, _ := cli.GetStream(vid, vid.Formats.WithAudioChannels().FindByQuality("720p"))

			fileName := execPath + "\\..\\videos\\" + strings.Replace(result.Title, " ", "-", -1) + ".mp4"
			// Make a file to write the video to.
			f, _ := os.Create(fileName)
			// Make a reader with a cancelable context.
			readerctx := NewReader(ctx, reader)

			// Copy the data from the video stream to the file.
			_, err := io.Copy(f, readerctx)
			if err != nil {
				os.Remove(fileName)
			}
			f.Close()
		}(result, m.items[i])
	}

	m.PublishRowsReset()
}

func (m *DownloadModel) CancelAllDownloads() {
	for _, item := range m.items {
		item.Downloading = false
	}

	m.items = make([]*Download, 0)
	m.PublishRowsReset()
}
