package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// グローバル変数の宣言！（初期化はmain関数内で行う）
var discordToken string
var textChannelID string
var err error

func main() {
	// 環境変数のロード
	loadEnv()
	handler()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// トークンを環境変数から取得
	discordToken = os.Getenv("DISCORDTOKEN")
	if discordToken == "" {
		log.Fatal("DISCORDTOKEN is not set in .env")
	}
}

func handler() {
	// Discord APIに接続
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Discordセッションの作成に失敗しました: %v", err)
	}

	// Botを起動し、Discordサーバーに接続
	err = dg.Open()
	if err != nil {
		log.Fatalf("Discordサーバーへの接続に失敗しました: %v", err)
	}

		// Discordセッションの使用後、自動的にクローズ
	defer dg.Close()

	// `Ctrl + C` で安全に終了できるようにする
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("Bot is running... Press Ctrl+C to exit.")

	// `stop` チャネルが値を受け取るまで待機
	<-stop

	log.Println("Bot is shutting down...")
}
