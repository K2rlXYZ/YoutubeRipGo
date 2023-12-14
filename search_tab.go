package main

import (
	"reflect"
	"regexp"

	youtube "github.com/KarlMul/youtubeGo"
	"github.com/lxn/walk"
)

func SearchQuery(query string) (bool, any) {
	var client = youtube.Client{
		YtApiv3Key: YtApiv3KEY,
	}

	var videoPattern, err = regexp.Compile(`(?:watch\?v=)([\w-]*)|(?:https:\/\/youtu\.be\/)([\w-]*)(?:\?si)`)
	if err != nil {
		walk.MsgBox(nil, "Error", "Unable to compile regex, "+err.Error(), walk.MsgBoxOK)
	}
	var matchedVideoUrl = videoPattern.FindString(query)
	if matchedVideoUrl != "" {
		var video, err = client.GetVideo(query)
		if err != nil {
			walk.MsgBox(nil, "Error", "Unable to search query, video search error:\n"+err.Error(), walk.MsgBoxOK)
			return false, nil
		}
		return true, video
	}

	var playlistPattern, errr = regexp.Compile(`(?:playlist\?list=)([\w-]*)(?:&si|$)`)
	if errr != nil {
		walk.MsgBox(nil, "Error", "Unable to search query, playlist search error:\n"+err.Error(), walk.MsgBoxOK)
		return false, nil
	}
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

func OnSearch() {
	var found, res = SearchQuery(inSearchQuery.Text())
	switch reflect.TypeOf(res) {
	case reflect.TypeOf(&youtube.Video{}):
		resultsTableModel.SetResultRowsFromVideo(res.(*youtube.Video))
	case reflect.TypeOf(&youtube.QueryResponseData{}):
		resultsTableModel.SetResultRowsFromQueryResponseData(res.(*youtube.QueryResponseData))
	case reflect.TypeOf([]*youtube.PlaylistEntry{}):
		resultsTableModel.SetResultRowsFromVideos(res.([]*youtube.PlaylistEntry))
	}
	if found {
		tabWidget.SetCurrentIndex(resultsTabEnum)
	}
}
