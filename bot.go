package albumbot

import (
	"fmt"
	"os"
	"os/signal"
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
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
	}

	if m.Content == "!taisho" {
		urls, e := GetAlbumUrls("taisho")
		//var errortxt = e
		fmt.Println(e)
		//s.ChannelMessageSend(m.ChannelID, printError(e))
		s.ChannelMessageSend(m.ChannelID, urls[0])
	}

	if m.Content == "!oemori" {
		s.ChannelMessageSend(m.ChannelID, "oemori")
	}

	if m.Content == "!test" {
		titles, err := GetAlbumTitles()
		if err != nil {
			fmt.Println(err)
		}
		s.ChannelMessageSend(m.ChannelID, strings.Join(titles, ","))
	}

	/*if strings.Contains(m.Content, "title:") && strings.Contains(m.Content, "urls:") {
		var tmp = m.ContentWithMentionsReplaced()
	}*/
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
