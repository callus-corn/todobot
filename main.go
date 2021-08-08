package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	TOKEN = ""
)

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	//メモ追加
	dg.AddHandler(add)
	//メモ削除
	dg.AddHandler(remove)
	//ロット
	dg.AddHandler(loot)
	//ロットバトル
	dg.AddHandler(lootBattle)
	//周回ロット
	dg.AddHandler(lootGrinding)
	//出るまでロット
	dg.AddHandler(lootUntilGet)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	//Ctrl+Cで終了
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
	/*
	   L:
	   	for {
	   		select {
	   		case <-timer.C:
	   			rotate()
	   			check(dg, "")
	   			reset, err := getTime()
	   			if err != nil {
	   				log.Fatal(err)
	   			}
	   			timer.Reset(time.Until(reset))
	   		case <-sc:
	   			dg.Close()
	   			break L
	   		}
	   	}
	*/
}

/*
func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "私の祖父ルイゾワも、第七霊災を防ぐため、カルテノーで、エオルゼア十二神の神降ろしを試みている。･･････だからイゼル、あなたの気持ちもわからないではない。しかし、グナース族の望みは単純な領土欲だ･･････。もしそれが本当だとすれば、あまりに無邪気に思える。やはり、彼らの蛮神を討伐しなければなるまい。" {
		s.ChannelMessageSend(m.ChannelID, "言うは易しだな、アルフィノ……。お前がグナース族の蛮神と戦うというのなら別だが、蛮神討伐となれば、「光の戦士」に頼るほかあるまい？")
	}

}
*/

func add(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	query := strings.Split(m.Content, "/")

	if len(query) == 3 && query[0] == "追加" {
		values := url.Values{}
		values.Set("target", query[1])
		values.Add("task", query[2])

		res, err := http.PostForm("http://18.217.133.253/add", values)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		s.ChannelMessageSend(m.ChannelID, "http://18.217.133.253/")
	}
}

func remove(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	query := strings.Split(m.Content, "/")

	if len(query) == 3 && query[0] == "削除" {
		values := url.Values{}
		values.Set("target", query[1])
		values.Add("task", query[2])

		res, err := http.PostForm("http://18.217.133.253/remove", values)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		s.ChannelMessageSend(m.ChannelID, "http://18.217.133.253/")
	}
}

func loot(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	username := m.Author.Username
	switch m.Author.Username {
	case "あおさ":
		username = "こせき"
	case "たこのめ":
		username = "みつ"
	case "willow":
		username = "ヤギヌマ"
	case "氷筍":
		username = "タケシ"
	case "salt_rippi":
		username = "そると"
	}

	if m.Content == "ロット" {
		result := rand.Intn(99) + 1
		s.ChannelMessageSend(m.ChannelID, username+"はダイスで"+strconv.Itoa(result)+"を出した。")
	}
}

func lootBattle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	query := strings.Split(m.Content, "/")

	if len(query) > 1 && query[0] == "ロットバトル" {
		dice := make(map[string]int)
		soldier := query[1:]
		max := 0
		winner := ""

		for _, v := range soldier {
			dice[v] = rand.Intn(99) + 1
			s.ChannelMessageSend(m.ChannelID, v+"はダイスで"+strconv.Itoa(dice[v])+"を出した。")
			if dice[v] > max {
				max = dice[v]
				winner = v
			}
		}

		s.ChannelMessageSend(m.ChannelID, winner+"は勝利を手に入れた。")
	}
}

func lootGrinding(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	username := m.Author.Username
	switch m.Author.Username {
	case "あおさ":
		username = "こせき"
	case "たこのめ":
		username = "みつ"
	case "willow":
		username = "ヤギヌマ"
	case "氷筍":
		username = "タケシ"
	case "salt_rippi":
		username = "そると"
	}

	query := strings.Split(m.Content, "/")

	if len(query) == 2 && query[0] == "周回ロット" {
		loop, err := strconv.Atoi(query[1])
		if err != nil || loop < 1 {
			s.ChannelMessageSend(m.ChannelID, "使い方間違ってる")
			return
		}

		if loop > 99 {
			s.ChannelMessageSend(m.ChannelID, "アイテム交換しろ")
			return
		}

		soldier := []string{username, "敵1", "敵2", "敵3", "敵4", "敵5", "敵6", "敵7"}
		results := []string{}
		for i := 0; i < loop; i++ {
			max := 0
			winner := ""
			line := ""
			for _, v := range soldier {
				dice := rand.Intn(99) + 1
				if dice > max {
					max = dice
					winner = v
				}
				line += v + "：" + strconv.Itoa(dice) + "　"
			}
			if winner == username {
				line = "☆　" + line
				winner = "☆　" + winner
			}
			results = append(results, line)
			results = append(results, winner+"は勝利を手に入れた。")
		}

		result := strings.Join(results, "\n")
		s.ChannelFileSend(m.ChannelID, "result.txt", strings.NewReader(result))
	}
}

func lootUntilGet(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	username := m.Author.Username
	switch m.Author.Username {
	case "あおさ":
		username = "こせき"
	case "たこのめ":
		username = "みつ"
	case "willow":
		username = "ヤギヌマ"
	case "氷筍":
		username = "タケシ"
	case "salt_rippi":
		username = "そると"
	}

	if m.Content == "出るまでロット" {
		soldier := []string{username, "敵1", "敵2", "敵3", "敵4", "敵5", "敵6", "敵7"}
		winner := ""

		i := 0
		results := []string{}
		for winner != username {
			max := 0
			winner = ""
			line := ""
			for _, v := range soldier {
				dice := rand.Intn(99) + 1
				if dice > max {
					max = dice
					winner = v
				}
				line += v + "：" + strconv.Itoa(dice) + "　"
			}
			result := line + "\n" + winner + "は勝利を手に入れた。"
			results = append(results, result)
			i += 1
		}
		win := username + "は" + strconv.Itoa(i) + "回目に勝利を手に入れた。\n"
		out := strings.Join(results, "\n")
		s.ChannelFileSendWithMessage(m.ChannelID, win, "result.txt", strings.NewReader(out))
	}
}

/*
func ready(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	query := strings.Split(m.Content, "/")

	if len(query) == 2 && query[0] == "レディチェ" {
		file, err := os.OpenFile("ready", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		username := m.Author.Username
		switch m.Author.Username {
		case "あおさ":
			username = "こせき"
		case "たこのめ":
			username = "みつ"
		case "willow":
			username = "ヤギヌマ"
		case "氷筍":
			username = "タケシ"
		case "salt_rippi":
			username = "そると"
		}

		fmt.Fprintln(file, username+"　"+query[1])
		s.ChannelMessageSend(m.ChannelID, username+"のレディチェを受け付けました")
	}
}

func check(s *discordgo.Session, c string) {
	file, err := os.OpenFile("ready", os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ready, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	s.ChannelMessageSend(c, string(ready))

	ioutil.WriteFile("ready", []byte("今日やる感じ？やらない感じ？\r\n"), 0600)
}

func getTime() (time.Time, error) {
	file, err := os.OpenFile("date", os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	date, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	d := strings.Split(string(date), "\r\n")[0]

	return time.Parse("2006-01-02 15:04:05  (MST)", d)
}

func rotate() {
	file, err := os.OpenFile("date", os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	date := strings.Split(string(data), "\r\n")

	old, err := getTime()
	if err != nil {
		log.Fatal(err)
	}

	new := append(date[1:], old.AddDate(0, 0, 7).Format("2006-01-02 15:04:05  (MST)"))
	//new := append(date[1:], old.Add(time.Second*30).Format("2006-01-02 15:04:05  (MST)"))
	result := strings.Join(new, "\r\n")

	ioutil.WriteFile("date", []byte(result), 0600)

}
*/
