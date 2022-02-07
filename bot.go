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

func New() {
	discordToken := "Bot " + os.Getenv("DISCORD_TOKEN")

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
	for i := range s {
		if s == arr[i] {
			return i, true
		}
	}
	return 0, false
}

func albumadd(s *discordgo.Session, m *discordgo.MessageCreate) error {
	contents := strings.Split(m.Content, " ")
	if len(contents) != 2 {
		return fmt.Errorf("→ !albumadd actual_albumname の形でファイルをアップロードしてね！")
	}
	title := contents[1]
	titles, err := GetAlbumTitles()
	if err != nil {
		return err
	}
	if !contains(titles, title) {
		return fmt.Errorf("%sというアルバムはなかったよ。!albumcreateコマンドで作れるよ！", title)
	}
	if len(m.Attachments) == 0 {
		return fmt.Errorf("画像が一枚も添付されてないよ。")
	}
	invalidAttaches := []string{}
	for _, attach := range m.Attachments {
		if isUrlImage(attach.URL) {
			err := PostAlbumUrl(title, attach.URL)
			s.ChannelMessageSend(m.ChannelID, attach.URL+" を"+title+"アルバムに追加したよ！")
			if err != nil {
				return err
			}
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
	return `・登録されているアルバムから見たいアルバムを選択する
!album 
・アルバムを作成する
!albumcreate albumtitle
・アルバムに写真を追加する（以下のコマンドと同時に写真を添付）
!albumadd actual_albumname
	`
}
func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "!Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
	}
	if (m.Content == "!albumbot -h") || (m.Content == "!albumbot -help") || (m.Content == "!albumbot help") {
		s.ChannelMessageSend(m.ChannelID, checkclhelp())
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
		if len(titles) <= 10 {
			for i, v := range titles {
				s.ChannelMessageSend(m.ChannelID, getNumEmoji(i+1)+" "+v)
			}
			s.ChannelMessageSend(m.ChannelID, "番号を選んでね！")
		} else {
			for i := 0; i < 10; i++ {
				s.ChannelMessageSend(m.ChannelID, getNumEmoji(i+1)+" "+titles[i])
			}
			s.ChannelMessageSend(m.ChannelID, "番号を選んでね！")
		}
	}

	if strings.HasPrefix(m.Content, "!albumcreate") {
		arr1 := strings.Split(m.Content, " ")
		if len(arr1) == 2 && arr1[0] == "!albumcreate" {
			err := CreateAlbum(arr1[1])
			if err != nil {
				return err
			}
			s.ChannelMessageSend(m.ChannelID, arr1[1]+"というアルバムを作成したよ！")
		} else {
			s.ChannelMessageSend(m.ChannelID, "→ !albumcreate titlename の形で記入してね！")
		}
	}

	if strings.HasPrefix(m.Content, "!albumadd") {
		err := albumadd(s, m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
	}

	if m.Content == "番号を選んでね！" && m.Author.ID == s.State.User.ID {
		titles, err := GetAlbumTitles()
		if err != nil {
			panic(err)
		}
		if len(titles) <= 10 {
			for i := 0; i < len(titles); i++ {
				s.MessageReactionAdd(m.ChannelID, m.ID, getNumEmoji(i+1))
			}
		} else {
			for i := 0; i < 10; i++ {
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
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		panic(err)
	}
	//botが投稿した"番号を選んでね！"のメッセージのみ処理
	if r.UserID != s.State.User.ID && message.Content == "番号を選んでね！" && message.Author.ID == s.State.User.ID {
		if r.MessageReaction.Emoji.Name == "➡️" { //アルバムのページを進める操作予定

		} else if r.MessageReaction.Emoji.Name == "⬅️" { //アルバムのページを戻す操作予定

		}

		index, NumEmojiFlag := getNumFromNumEmoji(r.MessageReaction.Emoji.Name)
		if NumEmojiFlag {
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)

			urls, err := GetAlbumUrls(titles[index])
			if err != nil {
				panic(err)
			}
			s.ChannelMessageSend(r.ChannelID, "> "+titles[index])
			for _, url := range urls {
				s.ChannelMessageSend(r.ChannelID, url)
			}

		}
	} else {

	}

}
