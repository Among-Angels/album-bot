package albumbot

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func New() {
	discordToken := "Bot " + loadToken()

	session, err := discordgo.New()
	if err != nil {
		fmt.Println("Error in create session")
		panic(err)
	}
	session.Token = discordToken
	session.AddHandler(onMessageCreate)
	session.AddHandler(onReactionAdd)

	if err = session.Open(); err != nil {
		panic(err)
	}
	defer session.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	fmt.Println("booted!!!")

	<-sc
	return
}

//getter関数を定義
func getNumOptions() []string {
	arr := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣"}
	return arr
}

//数字から数字スタンプ文字列を返す
func getNumEmoji(i int) string {
	if i < 1 {
		return "❓"
	}
	// 対応する絵文字がない場合はその値をそのまま返す
	if i > 9 {
		return strconv.Itoa(i)
	}
	arr := getNumOptions()
	return arr[i-1]
}

//数字スタンプ文字列から数値とbool値を返す
func getNumFromNumEmoji(s string) (int, bool) {
	arr := getNumOptions()
	for i := range s {
		if s == arr[i] {
			return i, true
		}
	}
	return 0, false
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "!Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
	}

	if m.Content == "!taisho" {
		urls, e := GetAlbumUrls("taisho")
		fmt.Println(e)
		s.ChannelMessageSend(m.ChannelID, urls[0])
	}

	if m.Content == "!album" {
		titles, err := GetAlbumTitles()
		if err != nil {
			panic(err)
		}
		if len(titles) <= 9 {
			for i, v := range titles {
				s.ChannelMessageSend(m.ChannelID, strconv.Itoa(i+1)+"."+v)
			}
			s.ChannelMessageSend(m.ChannelID, "番号を選んでね！")
		} else {
			for i := 0; i < 9; i++ {
				s.ChannelMessageSend(m.ChannelID, strconv.Itoa(i+1)+"."+titles[i])
			}
			s.ChannelMessageSend(m.ChannelID, "番号を選んでね！")
		}
	}
	if strings.Contains(m.Content, "!createalbum") {
		arr1 := strings.Split(m.Content, "")
		if len(arr1) == 2 && arr1[0] == "!createalbum" {
			CreateAlbum(arr1[1])
		} else {
			s.ChannelMessageSend(m.ChannelID, "!albumcreate titlename の形で記入してね！")
		}
	}

	if m.Content == "番号を選んでね！" && m.Author.ID == s.State.User.ID {
		titles, err := GetAlbumTitles()
		if err != nil {
			panic(err)
		}
		if len(titles) <= 9 {
			for i := 0; i < len(titles); i++ {
				s.MessageReactionAdd(m.ChannelID, m.ID, getNumEmoji(i+1))
			}
		} else {
			for i := 0; i < 9; i++ {
				s.MessageReactionAdd(m.ChannelID, m.ID, getNumEmoji(i+1))
			}
			s.MessageReactionAdd(m.ChannelID, m.ID, "➡️")
		}
	}

}

func onReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	titles, err := GetAlbumTitles()
	if err != nil {
		panic(err)
	}
	if r.UserID != s.State.User.ID {
		if r.MessageReaction.Emoji.Name == "➡️" { //アルバムのページを進める操作予定

		} else if r.MessageReaction.Emoji.Name == "⬅️" { //アルバムのページを戻す操作予定

		}

		index, flag := getNumFromNumEmoji(r.MessageReaction.Emoji.Name)
		if flag {
			urls, err := GetAlbumUrls(titles[index])
			if err != nil {
				panic(err)
			}
			for _, url := range urls {
				s.ChannelMessageSend(r.ChannelID, url)
			}
			s.ChannelMessageSend(r.ChannelID, r.MessageReaction.Emoji.ID)
			flag = false
		}
	} else {

	}

}

func loadToken() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("cannot load envrionments: %v", err)
	}
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		panic("no discord token exists.")
	}
	return token
}
