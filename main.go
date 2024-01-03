// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

var projectID string
var bucketName string
var geminiKey string

func main() {
	var err error
	projectID = os.Getenv("GCS_PROJECT_ID")
	bucketName = os.Getenv("GCS_BUCKET_NAME")
	geminiKey = os.Getenv("GOOGLE_GEMINI_API_KEY")

	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	if err != nil {
		log.Println("Bot:", bot, " err:", err)
	}
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				req := message.Text

				// Reply with Gemini result
				ret := GeminiChatComplete(req)

				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret)).Do(); err != nil {
					log.Print(err)
				}
			// Handle only on Sticker message
			case *linebot.StickerMessage:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
					log.Print(err)
				}

			// Handle only image message
			case *linebot.ImageMessage:
				log.Println("Got img msg ID:", message.ID)

				//Get image binary from LINE server based on message ID.
				content, err := bot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Println("Got GetMessageContent err:", err)
				}
				defer content.Content.Close()
				data, err := io.ReadAll(content.Content)
				if err != nil {
					log.Fatal(err)
				}
				ret, err := GeminiImage(data)
				if err != nil {
					ret = "無法辨識影片內容文字，請重新輸入:" + err.Error()
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ret)).Do(); err != nil {
					log.Print(err)
				}

			// Handle only video message
			case *linebot.VideoMessage:
				log.Println("Got video msg ID:", message.ID)
			}
		}
	}
}
