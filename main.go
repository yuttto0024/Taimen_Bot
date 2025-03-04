package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

// グローバル変数の宣言！（初期化はmain関数内で行う）
var discordToken string
var textChannelID string

func main() {
	lambda.Start(handler)
	//handler()
}

func handler() {

	// 環境変数から取得
	discordToken = os.Getenv("DISCORDTOKEN")
	if discordToken == "" {
		log.Fatal("環境変数 DISCORDTOKEN が設定されていません")
	}

	textChannelID = os.Getenv("DISCORDTEXTCHANNELID")
	if textChannelID == "" {
		log.Fatal("環境変数 DISCORDTEXTCHANNELID が設定されていません")
	}
	
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

	// 上位3名の情報をEmbed/通常のメッセージ形式に組み立てて送信
	sendMessages(dg, textChannelID)

	defer dg.Close() // Lambdaなので必ずセッションを閉じる
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

