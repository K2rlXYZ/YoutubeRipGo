package main

import (
	"regexp"

	youtube "github.com/KarlMul/youtubeGo"
	"github.com/lxn/walk"
)

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
