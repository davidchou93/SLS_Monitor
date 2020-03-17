package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var (
	ENDPOINT   = os.Getenv("CRYPTOWATCH")
	OHLCSTRUCT = [7]string{"closeTime", "openPrice", "highPrice", "lowPrice", "closePrice", "volume", "quoteVolume"}
	BOT_TOKEN  = os.Getenv("BOT_TOKEN")
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	now := time.Now().Unix()
	after := time.Now().Add(-5 * time.Minute).Unix()
	result := map[string]interface{}{}
	todoList := []string{"btcusdt"}
	bot, _ := tgbotapi.NewBotAPI(BOT_TOKEN)
	msg := tgbotapi.NewMessageToChannel("-1001275593710", "")
	for _, s := range todoList {
		target := fmt.Sprintf("%s/markets/binance/%s/ohlc?periods=300&before=%d&after=%d", ENDPOINT, s, now, after)
		msg.Text += fmt.Sprintln("======")
		msg.Text += fmt.Sprintln(strings.ToUpper(s))
		r, err := http.Get(target)
		if err != nil {
			return Response{StatusCode: 404, Body: err.Error()}, nil
		}
		defer r.Body.Close()
		rBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err.Error())
		}

		rMap := map[string]map[string]interface{}{
			"result":    make(map[string]interface{}),
			"allowance": make(map[string]interface{}),
		}
		err = json.Unmarshal(rBody, &rMap)
		if err != nil {
			log.Printf("Unmarshal failed:%s", err.Error())
		}
		info := rMap["result"]["300"].([]interface{})[0]
		OHLC := map[string]float64{}
		for index, value := range OHLCSTRUCT {
			OHLC[value] = info.([]interface{})[index].(float64)
		}

		msg.Text += fmt.Sprintln("======")
		priceChange := (OHLC["closePrice"] - OHLC["openPrice"]) / OHLC["openPrice"]
		fluctuation := (OHLC["highPrice"] - OHLC["lowPrice"]) / OHLC["openPrice"]
		if priceChange >= 0.01 || fluctuation >= 0.05 {
			msg.Text += fmt.Sprintln("<- NEED NOTICE ->")
			msg.Text += fmt.Sprintln("======")
		}
		msg.Text += fmt.Sprintf("%s:%.2f%%\n", "PriceChange", priceChange*100)
		msg.Text += fmt.Sprintf("%s:%.2f%%\n", "Fluctuation", fluctuation*100)
		msg.Text += fmt.Sprintln("======")

		msg.Text += fmt.Sprintf("%s:  %.2f\n", "OpenPrice", OHLC["openPrice"])
		msg.Text += fmt.Sprintf("%s:  %.2f\n", "ClosePrice", OHLC["closePrice"])
		msg.Text += fmt.Sprintf("%s:  %.2f\n", "HighPrice", OHLC["highPrice"])
		msg.Text += fmt.Sprintf("%s:  %.2f\n", "LowPrice", OHLC["lowPrice"])
		msg.Text += fmt.Sprintf("%s:  %.1f\n", "Volume", OHLC["volume"])
		msg.Text += fmt.Sprintf("%s:  %.1f\n", "QuoteVolume", OHLC["quoteVolume"])
		msg.Text += fmt.Sprintf("%s:  %s\n", "CloseTime", time.Unix(int64(OHLC["closeTime"]), 0).Format("20060102-15:05"))

		result[s] = make(map[string]float64)
		for key, value := range OHLC {
			result[s].(map[string]float64)[key] = value
		}
		msg.Text += fmt.Sprintln("======")
	}
	_, sendErr := bot.Send(msg)
	if sendErr != nil {
		fmt.Println(sendErr.Error())
	}
	body, _ := json.Marshal(result)
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "echo-handler",
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
