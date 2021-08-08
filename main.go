package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	TOKEN = "ODcwOTIyNzcxODU1NzA4MTkw.YQT0CQ.8sCsa4bY0Y2zNqymJN4NQPiL3sM"
)

func main() {
	//セリフ用乱数設定
	rand.Seed(int64(time.Now().Second()))

	//ready機能の開始時間設定
	t, err := getTime()
	if err != nil {
		log.Fatal(err)
	}
	start := time.Until(t)
	timer := time.NewTimer(start)

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	//返事
	dg.AddHandler(message)
	//メモ追加
	dg.AddHandler(add)
	//メモ削除
	dg.AddHandler(remove)
	//レディチェ
	dg.AddHandler(ready)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	//Ctrl+Cで終了
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

L:
	for {
		select {
		case <-timer.C:
			rotate()
			check(dg, "868827776873009292")
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
}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "私の祖父ルイゾワも、第七霊災を防ぐため、カルテノーで、エオルゼア十二神の神降ろしを試みている。･･････だからイゼル、あなたの気持ちもわからないではない。しかし、グナース族の望みは単純な領土欲だ･･････。もしそれが本当だとすれば、あまりに無邪気に思える。やはり、彼らの蛮神を討伐しなければなるまい。" {
		s.ChannelMessageSend(m.ChannelID, "言うは易しだな、アルフィノ……。お前がグナース族の蛮神と戦うというのなら別だが、蛮神討伐となれば、「光の戦士」に頼るほかあるまい？")
	}

}

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

		random := []string{
			"言うは易しだな、" + query[1] + "……。お前が一人で戦うというのなら別だが、零式攻略となれば、「光の戦士」に頼るほかあるまい？",
			"薪拾いなら任せてくれよ！",
		}
		s.ChannelMessageSend(m.ChannelID, random[rand.Intn(len(random))])
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

		random := []string{
			"フフ……やはり、お前は……笑顔が……イイ……。",
			"さらばだ光の戦士　私を導いてくれてありがとう",
		}
		s.ChannelMessageSend(m.ChannelID, random[rand.Intn(len(random))])
	}
}

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
