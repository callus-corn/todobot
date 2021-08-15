package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
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

	//タイマー生成
	dg.AddHandler(create)
	//タイマーリスト
	dg.AddHandler(list)
	//タイマー開始
	dg.AddHandler(play)
	//タイマー停止
	dg.AddHandler(stop)
	//タイマー削除
	dg.AddHandler(delete)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	//Ctrl+Cで終了
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
