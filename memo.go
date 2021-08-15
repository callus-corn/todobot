package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
			s.ChannelMessageSend(m.ChannelID, "通信エラー")
			return
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
			s.ChannelMessageSend(m.ChannelID, "通信エラー")
			return
		}
		defer res.Body.Close()

		s.ChannelMessageSend(m.ChannelID, "http://18.217.133.253/")
	}
}
