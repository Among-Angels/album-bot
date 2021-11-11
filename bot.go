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
		for i, v := range titles {
			s.ChannelMessageSend(m.ChannelID, strconv.Itoa(i)+"."+v)
		}
		s.ChannelMessageSend(m.ChannelID, "番号を選んでね！")
	}

	if m.Content == "番号を選んでね！" && m.Author.ID == s.State.User.ID {
		s.MessageReactionAdd(m.ChannelID, m.ID, "1️⃣")
		s.MessageReactionAdd(m.ChannelID, m.ID, "2️⃣")
	}

}

func onReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	titles, err := GetAlbumTitles()
	if err != nil {
		panic(err)
	}
	if r.UserID != s.State.User.ID && r.MessageReaction.Emoji.Name == "1️⃣" {
		urls, err := GetAlbumUrls(titles[0])
		if err != nil {
			panic(err)
		}
		for _, url := range urls {
			s.ChannelMessageSend(r.ChannelID, url)
		}
		s.ChannelMessageSend(r.ChannelID, r.MessageReaction.Emoji.ID)
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
