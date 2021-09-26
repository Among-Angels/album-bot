package albumbot

import (
	"reflect"
	"testing"
)

func TestGetAlbumUrls(t *testing.T) {
	wants := []string{
		"https://cdn.discordapp.com/attachments/723903982745157723/747799894001189024/00100lPORTRAIT_00100_BURST20190525180748800_COVER.jpg",
		"https://cdn.discordapp.com/attachments/723903982745157723/782582360688033842/image0.jpg",
		"https://cdn.discordapp.com/attachments/723903982745157723/784664833860960266/PXL_20201128_013520140.jpg",
		"https://cdn.discordapp.com/attachments/723903982745157723/788395055843901450/image0.jpg",
		"https://cdn.discordapp.com/attachments/723903982745157723/853546288527441940/unknown.png",
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
			title: "taisho",
		},
		want: wants,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := GetAlbumUrls(tt.args.title); !reflect.DeepEqual(got, tt.want) {
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
