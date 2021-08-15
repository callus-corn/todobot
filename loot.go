package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
