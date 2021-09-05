package albumbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Albumはタイトルとそれに紐付けられた画像URLの集合です。
type Album struct {
	Title string
	Urls  []string
}

// Albumsはアルバムデータが保存されたJsonファイル全体を表現します。
type Albums struct {
	Albums []Album
}

// GetAlbumUrlsは与えられたアルバム名の画像のURLのリストを返します。
func GetAlbumUrls(title string) (urls []string, e error) {
	raw, err := ioutil.ReadFile("./dataSet.json")
	if err != nil {
		return nil, err
	}

	var albums Albums
	json.Unmarshal(raw, &albums)
	for _, a := range albums.Albums {
		if a.Title == title {
			return a.Urls, nil
		}
	}
	e = fmt.Errorf("アルバム%sが見つかりませんでした。", title)
	return nil, e
}
