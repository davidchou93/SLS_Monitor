package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var (
	BOT_TOKEN = os.Getenv("BOT_TOKEN")
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	bot, err := tgbotapi.NewBotAPI(BOT_TOKEN)
	bot.Debug = true
	update := tgbotapi.Update{}
	err = json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return Response{StatusCode: 400}, nil
	}
	if update.Message == nil {
		return Response{StatusCode: 200, Body: "Empty message from TG, do nothing."}, nil
	}
	if update.Message.IsCommand() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "echo":
			arguments := update.Message.CommandArguments()
			if strings.HasPrefix(arguments, "@") {
				msg.Text = strings.SplitN(arguments, " ", 2)[1]
			} else {
				msg.Text = arguments
			}
		case "speak":
			arguments := update.Message.CommandArguments()
			message := ""
			if strings.HasPrefix(arguments, "@") {
				message = strings.SplitN(arguments, " ", 2)[1]
			} else {
				message = arguments
			}
			if strings.Count(message, "")-1 >= 2 {
				svc := polly.New(session.New())
				input := &polly.SynthesizeSpeechInput{
					// LexiconNames: []*string{
					// aws.String(fmt.Sprintf("voice%d",time.Unix()),
					// },
					LanguageCode: aws.String("cmn-CN"),
					OutputFormat: aws.String("ogg_vorbis"),
					SampleRate:   aws.String("8000"),
					Text:         aws.String(message),
					TextType:     aws.String("text"),
					VoiceId:      aws.String("Zhiyu"),
				}
				v, err := svc.SynthesizeSpeech(input)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						log.Println(aerr.Error())
						return Response{StatusCode: 400, Body: "Polly failed to transfer message."}, nil
					} else {
						log.Println(err.Error())
					}
				}
				// Transfer audiofile into []Bytes
				audioFile, err := ioutil.ReadAll(v.AudioStream)
				fileBytes := tgbotapi.FileBytes{
					Name:  fmt.Sprintf("voice%d", time.Now().Unix()),
					Bytes: audioFile,
				}
				voiceConfig := tgbotapi.NewVoiceUpload(update.Message.Chat.ID, fileBytes)
				bot.Send(voiceConfig)
			}
			msg.Text = ""
		default:
			msg.Text = "I don't know that command"
		}
		if msg.Text != "" {
			bot.Send(msg)
		}
	}
	result, _ := json.Marshal(map[string]string{"message": "succeed"})
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(result),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "receiver-handler",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
