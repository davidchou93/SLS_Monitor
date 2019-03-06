package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var(
	BOT_TOKEN   = os.Getenv("BOT_TOKEN")
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
	if update.Message == nil{
		return Response{StatusCode: 200,Body:"Empty message from TG, do nothing."},nil
	}
	
	if update.Message.IsCommand() {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,"")
		switch update.Message.Command(){
		case "help":
			msg.Text = "type /sayhi or /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}
		bot.Send(msg)
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
