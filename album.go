package albumbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

// GetAlbumPageは与えられたアルバム名と開始index,数量からURLのリストを返します
func GetAlbumPage(title string, start, count int) (urls []string, e error) {

	if start < 0 {
		e = fmt.Errorf("startは0以上の数値を指定してください")
		return nil, e
	}

	if count < 1 {
		e = fmt.Errorf("countは1以上の数値を指定してください")
		return nil, e
	}

	allUrls, err := GetAlbumUrls(title)
	if err != nil {
		return nil, err
	}

	if len(allUrls)-1 < start {
		e = fmt.Errorf("startはアルバム内のURL数より小さい値を指定してください")
		return nil, e
	}

	// count込で溢れた場合末尾まで返す
	if len(allUrls) < start+count {
		return allUrls[start:], nil
	}

	return allUrls[start : start+count], nil
}

// PostAlbumUrlは与えられたアルバム名のUrl配列に与えられたUrlを追加します
// ファイル全部読んで全部上書きする脳筋処理なので改良の余地ありです。。
func PostAlbumUrl(albumTitle, url string) (e error) {
	raw, err := ioutil.ReadFile("./dataSet.json")
	if err != nil {
		return err
	}

	var albums Albums
	json.Unmarshal(raw, &albums)
	for index, album := range albums.Albums {
		if album.Title != albumTitle {
			continue
		}
		albums.Albums[index].Urls = append(albums.Albums[index].Urls, url)
		marshaled, err := json.Marshal(albums)
		if err != nil {
			return err
		}
		writeError := ioutil.WriteFile("dataSet.json", marshaled, os.ModePerm)
		if writeError != nil {
			return writeError
		}
		return
	}
	e = fmt.Errorf("アルバム%sが見つかりませんでした。", albumTitle)
	return e
}
