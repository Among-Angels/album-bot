package albumbot

import (
	"reflect"
	"testing"
)

func TestGetAlbumTitles(t *testing.T) {
	wants := []string{"taisho", "oemori", "blank", "test"}
	tests := []struct {
		name       string
		wantTitles []string
		wantErr    bool
	}{{
		name:       "test",
		wantTitles: wants,
		wantErr:    false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitles, err := GetAlbumTitles()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlbumTitles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTitles, tt.wantTitles) {
				t.Errorf("GetAlbumTitles() = %v, want %v", gotTitles, tt.wantTitles)
			}
		})
	}
}
func TestGetAlbumUrls(t *testing.T) {
	wants := []string{
		"https://test1.png",
		"https://test2.png",
	}
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{{
		name: "test",
		args: args{
			title: "test",
		},
		want: wants,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAlbumUrls(tt.args.title)
			if err != nil {
				panic(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlbumUrls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAlbumPage(t *testing.T) {
	normalWants := []string{
		"https://cdn.discordapp.com/attachments/723903982745157723/782582360688033842/image0.jpg",
		"https://cdn.discordapp.com/attachments/723903982745157723/784664833860960266/PXL_20201128_013520140.jpg",
		"https://cdn.discordapp.com/attachments/723903982745157723/788395055843901450/image0.jpg",
	}
	overCountWants := []string{
		"https://cdn.discordapp.com/attachments/723903982745157723/788395055843901450/image0.jpg",
		"https://cdn.discordapp.com/attachments/723903982745157723/853546288527441940/unknown.png",
	}
	type args struct {
		title string
		start int
		count int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "normal case",
			args: args{
				title: "taisho",
				start: 1,
				count: 3,
			},
			want: normalWants,
		},
		{
			name: "over count case",
			args: args{
				title: "taisho",
				start: 3,
				count: 9999,
			},
			want: overCountWants,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := GetAlbumPage(tt.args.title, tt.args.start, tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlbumPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostAlbumUrl(t *testing.T) {
	type args struct {
		albumTitle string
		url        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "error",
			args: args{
				albumTitle: "invisible taisho",
				url:        "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostAlbumUrl(tt.args.albumTitle, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("PostAlbumUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
