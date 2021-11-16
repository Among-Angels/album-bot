package albumbot

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
func getNumOptions() []string {
	arr := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣"}
	return arr
}
func containAtIndexNum(s []string, e string) (int, bool) {
	for i, v := range s {
		if e == v {
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

	if m.Content == "番号を選んでね！" && m.Author.ID == s.State.User.ID {
		options := getNumOptions()
		titles, err := GetAlbumTitles()
		if err != nil {
			panic(err)
		}
		if len(titles) <= 9 {
			for i := 0; i < len(titles); i++ {
				s.MessageReactionAdd(m.ChannelID, m.ID, options[i])
			}
		} else {
			for i := 0; i < 9; i++ {
				s.MessageReactionAdd(m.ChannelID, m.ID, options[i])
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
		options := getNumOptions()
		index, flag := containAtIndexNum(options, r.MessageReaction.Emoji.Name)
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
