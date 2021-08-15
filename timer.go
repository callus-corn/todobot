package main

import (
	"context"
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	tts "cloud.google.com/go/texttospeech/apiv1"
	"github.com/bwmarrin/discordgo"
	ttspb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"gopkg.in/hraban/opus.v2"
)

func play(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	query := strings.Split(m.Content, "/")

	if len(query) == 3 && query[0] == "タイマー再生" {
		name := "タイマー/" + query[1]
		file, err := os.OpenFile(name+".txt", os.O_RDONLY, 0644)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマーが見つかりません")
			return
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマーデータ読み込みエラー")
			return
		}
		str := strings.NewReplacer("\r\n", "\n", "\r", "\n", "\n", "\n").Replace(string(bytes))
		lines := strings.Split(str, "\n")
		start, err := strconv.Atoi(query[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマーデータ解析エラー")
			return
		}

		go func() {
			timer := time.NewTimer(time.Duration(start) * time.Second)
			s.ChannelMessageSend(m.ChannelID, "タイマー開始まで"+strconv.Itoa(start)+"秒！")
			<-timer.C
			s.ChannelMessageSend(m.ChannelID, "タイマー開始！")
			now := time.Now()

			vc, err := s.ChannelVoiceJoin(m.GuildID, "868827776873009293", false, true)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "ごめん無理")
				return
			}
			defer vc.Disconnect()

			for _, line := range lines {
				t := strings.Split(line, "/")[0]
				min, err := strconv.Atoi(strings.Split(t, "-")[0])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "タイマー解析エラー")
					return
				}
				sec, err := strconv.Atoi(strings.Split(t, "-")[1])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "タイマー解析エラー")
					return
				}
				alertTime := now.Add(time.Duration(60*min+sec) * time.Second)
				timer.Reset(time.Until(alertTime))

				wave, err := ioutil.ReadFile(name + "/" + t + ".wav")
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "音声データエラー")
					return
				}
				for !reflect.DeepEqual(wave[:4], []byte{0x64, 0x61, 0x74, 0x61}) {
					wave = wave[1:]
				}
				size := int(binary.LittleEndian.Uint32(wave[4:8]))
				wave = wave[8:]

				source := make([][]byte, 0)

				for bytes := 0; bytes < size; bytes += 960 {
					//24kHz,20ms,1channel
					raw20ms := make([]int16, 480)
					for i := 0; i < 480; i++ {
						if len(wave) >= 2 {
							raw20ms[i] = int16(binary.LittleEndian.Uint16(wave[:2]))
							wave = wave[2:]
						} else {
							raw20ms[i] = 0
						}
					}
					audio := make([]byte, 1000)
					enc, err := opus.NewEncoder(24000, 1, opus.AppVoIP)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "エンコードエラー")
						return
					}
					n, err := enc.Encode(raw20ms, audio)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "音声フォーマットエラー")
						return
					}
					source = append(source, audio[:n])
				}

				<-timer.C
				if vc.Ready {
					vc.Speaking(true)
					for _, v := range source {
						vc.OpusSend <- v
					}
					vc.Speaking(false)
				} else {
					return
				}
			}
			timer.Reset(time.Until(time.Now().Add(time.Second)))
			<-timer.C
			s.ChannelMessageSend(m.ChannelID, "タイマー終了！")
		}()
	}
}

func stop(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "タイマー終了" {
		vc, err := s.ChannelVoiceJoin(m.GuildID, "868827776873009293", false, true)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "止められない！")
			return
		}
		vc.Disconnect()
	}
}

func list(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "タイマーリスト" {
		files, err := ioutil.ReadDir("タイマー")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "エラー")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "タイマー一覧")
		for _, file := range files {
			if strings.Contains(file.Name(), ".txt") {
				out := strings.NewReplacer(".txt", "").Replace(file.Name())
				s.ChannelMessageSend(m.ChannelID, out)
			}
		}
	}
}

func create(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	lines := strings.Split(m.Content, "\n")
	query := strings.Split(lines[0], "/")
	name := ""

	if len(query) == 2 && query[0] == "タイマー生成" {
		name = "タイマー/" + query[1]
		if err := os.Mkdir(name, 0644); err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマー作成エラー：01")
			return
		}
	} else {
		return
	}

	lines = lines[1:]
	if err := ioutil.WriteFile(name+".txt", []byte(strings.Join(lines, "\n")), 0644); err != nil {
		s.ChannelMessageSend(m.ChannelID, "タイマー作成エラー：02")
		return
	}

	texts := make(map[string]string)
	for _, v := range lines {
		time := strings.Split(v, "/")[0]
		text := strings.Split(v, "/")[1]
		texts[time] = text
	}

	for key, v := range texts {
		ctx := context.Background()

		client, err := tts.NewClient(ctx)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマー作成エラー：03")
			return
		}
		defer client.Close()

		req := ttspb.SynthesizeSpeechRequest{
			Input: &ttspb.SynthesisInput{
				InputSource: &ttspb.SynthesisInput_Text{Text: v},
			},

			Voice: &ttspb.VoiceSelectionParams{
				LanguageCode: "ja-JP",
				SsmlGender:   ttspb.SsmlVoiceGender_NEUTRAL,
			},
			AudioConfig: &ttspb.AudioConfig{
				AudioEncoding: ttspb.AudioEncoding_MP3,
			},
		}

		resp, err := client.SynthesizeSpeech(ctx, &req)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマー作成エラー：04")
			return
		}

		filename := name + "/" + key + ".mp3"
		err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマー作成エラー：05")
			log.Fatal(err)
			return
		}

		target := name + "/" + key
		if err := exec.Command("ffmpeg", "-i", target+".mp3", target+".wav").Run(); err != nil {
			s.ChannelMessageSend(m.ChannelID, "タイマー作成エラー：06")
			return
		}
	}
	s.ChannelMessageSend(m.ChannelID, "タイマーを作成しました")
}

func delete(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	query := strings.Split(m.Content, "/")

	if len(query) == 2 && query[0] == "タイマー削除" {
		if strings.Contains(query[1], "..") {
			s.ChannelMessageSend(m.ChannelID, "なんかヤバいことしようとしてない？")
			return
		}
		if err := os.RemoveAll("タイマー/" + query[1]); err != nil {
			s.ChannelMessageSend(m.ChannelID, "ディレクトリエラー")
			return
		}
		if err := os.Remove("タイマー/" + query[1] + ".txt"); err != nil {
			s.ChannelMessageSend(m.ChannelID, "ファイルエラー")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "タイマーを削除しました")
	}
}
