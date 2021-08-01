package main

import (
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
	TOKEN = ""
)

func main() {
	rand.Seed(int64(time.Now().Second()))

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(message)
	dg.AddHandler(add)
	dg.AddHandler(remove)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dg.Close()
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
		s.ChannelMessageSend(m.ChannelID, random[rand.Intn(2)])
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
		s.ChannelMessageSend(m.ChannelID, random[rand.Intn(2)])

	}
}
