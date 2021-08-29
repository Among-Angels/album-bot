package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
)

func main() {
    session, err := discordgo.New()
    if err != nil {
        fmt.Println("Error in create session")
        panic(err)
    }

    discordToken := loadToken()
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

func onMessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
    if event.Author.Bot {
        return
    }
    /* メッセージを受け取った際の処理 */
}
//DISCORD_TOKENは.envで管理予定
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