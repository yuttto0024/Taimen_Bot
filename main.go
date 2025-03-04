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

	textChannelID = os.Getenv("DISCORDTEXTCHANNELID")
	if textChannelID == "" {
		log.Fatal("TEXTCHANNELID is not set in .env")
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

	// 上位3名の情報をEmbed/通常のメッセージ形式に組み立てて送信
	sendMessages(dg, textChannelID)

	// `Ctrl + C` で安全に終了できるようにする
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("Bot is running... Press Ctrl+C to exit.")

	// `stop` チャネルが値を受け取るまで待機
	<-stop

	log.Println("Bot is shutting down...")
}

// メッセージを送信する関数
func sendMessages(s *discordgo.Session, channelID string) {
	message := "@everyone\n明日の対面活動に参加する人は <:sanka:1341527249236787302> 、欠席する人は <:kesseki:1342322128975954041> を押してね！！！"

	// メッセージを送信
	msg, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("メッセージの送信に失敗しました: %v", err)
		return
	}

	// 絵文字リアクションを追加
	reactions := []string{
		"sanka:1341527249236787302",
		"kesseki:1342322128975954041",
		"okuremasu:1342323653576101962",
		"online:1346094084812443749",
	}

	for _, reaction := range reactions {
		err = s.MessageReactionAdd(channelID, msg.ID, reaction) // `msg.ID` を適切に参照
		if err != nil {
			log.Printf("リアクションの追加に失敗しました: %v", err)
		}
	}
}

