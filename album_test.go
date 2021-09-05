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
