package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	TOKEN = ""
)

func main() {
	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(message)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if m.Content == "test" {
		s.ChannelMessageSend(m.ChannelID, "ok")
	}

	if m.Content == "私の祖父ルイゾワも、第七霊災を防ぐため、カルテノーで、エオルゼア十二神の神降ろしを試みている。･･････だからイゼル、あなたの気持ちもわからないではない。しかし、グナース族の望みは単純な領土欲だ･･････。もしそれが本当だとすれば、あまりに無邪気に思える。やはり、彼らの蛮神を討伐しなければなるまい。" {
		s.ChannelMessageSend(m.ChannelID, "言うは易しだな、アルフィノ……。お前がグナース族の蛮神と戦うというのなら別だが、蛮神討伐となれば、「光の戦士」に頼るほかあるまい？")
	}

	if m.Content == "追加 みつ マウントほしい" {
		json := `{"target":みつ,task="マウントほしい"}`
		req, err := http.NewRequest(
			"POST",
			"http://18.217.133.253",
			bytes.NewBuffer([]byte(json)),
		)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()

		s.ChannelMessageSend(m.ChannelID, "みつにメモを追加しました")
	}
}
