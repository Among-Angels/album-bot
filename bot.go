package albumbot

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var callCommand string

// New()の中で上書きされる可能性がある
// var titleindex = 0 //titleのindexを保持
// var pageindex = 0  //ページのindexを保持

var table = "Albums"

var currentBot *albumBot = &albumBot{}

type albumBot struct {
	//channelID        string
	//messageID        string
	//reactionID       string
	//albumSelectionID string
	urls []string
	//selectedTitle    string
	pageindex  int
	titleindex int
}

func newBot() *albumBot {
	return &albumBot{}
}

/* func (bot *albumBot) setAlbumSelectionID(id string) error {
	bot.albumSelectionID = id
} */

func New() {
	table = os.Getenv("TABLE_NAME")

	discordToken := "Bot " + os.Getenv("DISCORD_TOKEN")
	var ok bool
	callCommand, ok = os.LookupEnv("CALL_COMMAND")
	if !ok {
		callCommand = "!album"
	}
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)

	fmt.Println("booted!!!")

	<-sc
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func isUrlImage(url string) bool {
	exts := []string{"png", "jpg", "jpeg", "gif"}
	parts := strings.Split(url, ".")
	return contains(exts, parts[len(parts)-1])
}

//getter関数を定義
func getNumOptions() []string {
	arr := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}
	return arr
}

//数字から数字スタンプ文字列を返す
func getNumEmoji(i int) string {
	if i < 1 {
		return "❓"
	}
	// 対応する絵文字がない場合はその値をそのまま返す
	if i > 10 {
		return strconv.Itoa(i)
	}
	arr := getNumOptions()
	return arr[i-1]
}

//数字スタンプ文字列から数値とbool値を返す
func getNumFromNumEmoji(s string) (int, bool) {
	arr := getNumOptions()
	for i := range arr {
		if s == arr[i] {
			return i, true
		}
	}
	return 0, false
}

func albumAdd(s *discordgo.Session, m *discordgo.MessageCreate) error {
	contents := strings.Split(m.Content, " ")
	if len(contents) != 3 {
		return fmt.Errorf("→ " + callCommand + " add actual_albumname の形でファイルをアップロードしてね！")
	}
	title := contents[2]
	titles, err := GetAlbumTitles(table)
	if err != nil {
		return err
	}
	if !contains(titles, title) {
		return fmt.Errorf("%sというアルバムはなかったよ。"+callCommand+" createコマンドで作れるよ！", title)
	}
	if len(m.Attachments) == 0 {
		return fmt.Errorf("画像が一枚も添付されてないよ。")
	}
	invalidAttaches := []string{}
	for _, attach := range m.Attachments {
		if isUrlImage(attach.URL) {
			err := PostImage(table, title, attach.URL)
			if err != nil {
				return err
			}
			s.ChannelMessageSend(m.ChannelID, attach.URL+" を"+title+"アルバムに追加したよ！")
		} else {
			invalidAttaches = append(invalidAttaches, attach.Filename)
		}
	}

	if len(invalidAttaches) > 0 {
		return fmt.Errorf("以下のファイルは画像じゃないから無視したよ：\n%s", strings.Join(invalidAttaches, "\n"))
	}
	return nil
}

func checkclhelp() string {
	return callCommand + "\n・登録されているアルバムから見たいアルバムを選択する\n" +
		callCommand + " create albumtitle\n・アルバムを作成する\n" +
		callCommand + " add actual_albumname\n・アルバムに写真を追加する（以下のコマンドと同時に写真を添付）\n"

}

func commandSplit(str string) []string {
	commandArray := strings.Split(str, " ")
	return commandArray
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := commandSplit(m.Content)
	if len(command) == 0 || command[0] != callCommand {
		return
	}

	if len(command) == 1 {
		currentBot = newBot()

		titles, err := GetAlbumTitles(table)
		tmpstr := ""
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}

		if len(titles) > 10 {
			titles = titles[:10]
		}
		for i, v := range titles {
			tmpstr += getNumEmoji(i+1) + " " + v + "\n"
		}
		s.ChannelMessageSend(m.ChannelID, tmpstr)
		sent, err := s.ChannelMessageSend(m.ChannelID, "番号を選んでね！")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
		for i := range titles {
			s.MessageReactionAdd(m.ChannelID, sent.ID, getNumEmoji(i+1))
		}
		return
	}

	subCommand := command[1]
	switch subCommand {
	case "-h", "--help", "help":
		s.ChannelMessageSend(m.ChannelID, checkclhelp())
	case "create":
		if len(command) == 3 {
			err := CreateAlbum(table, command[2])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			}
			s.ChannelMessageSend(m.ChannelID, command[2]+"というアルバムを作成したよ！")
		} else {
			s.ChannelMessageSend(m.ChannelID, "→ "+callCommand+" create titlename の形で記入してね！")
		}
	case "add":
		err := albumAdd(s, m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
	}
}
func onReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	titles, err := GetAlbumTitles(table)
	tmpurl := ""
	if err != nil {
		s.ChannelMessageSend(r.ChannelID, err.Error())
	}
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		s.ChannelMessageSend(r.ChannelID, err.Error())
	}
	if r.UserID != s.State.User.ID && message.Content == "番号を選んでね！" && message.Author.ID == s.State.User.ID {
		index, NumEmojiFlag := getNumFromNumEmoji(r.MessageReaction.Emoji.Name)
		currentBot.titleindex = index
		if NumEmojiFlag {
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			urls, err := GetAlbumUrls(table, titles[index])
			if err != nil {
				s.ChannelMessageSend(r.ChannelID, err.Error())
			}
			currentBot.urls = urls
			s.ChannelMessageSend(r.ChannelID, "> "+titles[index])
			if len(urls) < 5 {
				for i := 0; i < len(urls); i++ {
					tmpurl += " " + urls[i]
				}
				s.ChannelMessageSend(r.ChannelID, tmpurl)
			} else {
				for i := 0; i < 5; i++ {
					tmpurl += " " + urls[i]
				}
				sent, err := s.ChannelMessageSend(r.ChannelID, tmpurl)
				if err != nil {
					s.ChannelMessageSend(r.ChannelID, err.Error())
				}
				s.MessageReactionAdd(r.ChannelID, sent.ID, "➡️")
			}
		}

		// 矢印押した人がbot以外である
		// メッセージのauthorがalbum botである
		// メッセージIDをグローバルで持つ必要がある？
	} else if r.UserID != s.State.User.ID && message.Author.ID == s.State.User.ID {
		if r.MessageReaction.Emoji.Name == "➡️" { //アルバムのページを進める操作
			currentBot.pageindex += 1
			if (5*currentBot.pageindex < len(currentBot.urls)) && (len(currentBot.urls) <= 5*(currentBot.pageindex+1)) { //pageindexが右端ページにあるときの処理
				for i := 0; i+5*currentBot.pageindex < len(currentBot.urls); i++ {
					tmpurl += " " + currentBot.urls[i+5*currentBot.pageindex]
				}
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				sent, err := s.ChannelMessageSend(r.ChannelID, tmpurl)
				if err != nil {
					s.ChannelMessageSend(r.ChannelID, err.Error())
				}
				s.MessageReactionAdd(r.ChannelID, sent.ID, "⬅")
			} else { //pageindexが中間ページにあるときの処理
				for i := 0; i+5*currentBot.pageindex < 5*(currentBot.pageindex+1); i++ {
					tmpurl += " " + currentBot.urls[i+5*currentBot.pageindex]
				}
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				sent, err := s.ChannelMessageSend(r.ChannelID, tmpurl)
				if err != nil {
					s.ChannelMessageSend(r.ChannelID, err.Error())
				}
				s.MessageReactionAdd(r.ChannelID, sent.ID, "⬅")
				s.MessageReactionAdd(r.ChannelID, sent.ID, "➡️")
			}
		} else if r.MessageReaction.Emoji.Name == "⬅" { //アルバムのページを戻す操作予定
			currentBot.pageindex -= 1
			if currentBot.pageindex == 0 { //pageindexが左端ページにあるときの処理
				for i := 0; i < 5; i++ {
					tmpurl += " " + currentBot.urls[i]
				}
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				sent, err := s.ChannelMessageSend(r.ChannelID, tmpurl)
				if err != nil {
					s.ChannelMessageSend(r.ChannelID, err.Error())
				}
				s.MessageReactionAdd(r.ChannelID, sent.ID, "➡️")
			} else { //pageindexが中間ページにあるときの処理
				for i := 0; i+5*currentBot.pageindex < len(currentBot.urls); i++ {
					tmpurl += " " + currentBot.urls[i+5*currentBot.pageindex]
				}
				s.ChannelMessageDelete(r.ChannelID, r.MessageID)
				sent, err := s.ChannelMessageSend(r.ChannelID, tmpurl)
				if err != nil {
					s.ChannelMessageSend(r.ChannelID, err.Error())
				}
				s.MessageReactionAdd(r.ChannelID, sent.ID, "⬅")
				s.MessageReactionAdd(r.ChannelID, sent.ID, "➡️")
			}
		}
	}
}
