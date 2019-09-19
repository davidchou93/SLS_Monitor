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
	OHLCSTRUCT = [6]string{"closeTime", "openPrice", "highPrice", "lowPrice", "closePrice", "volume"}
	BOT_TOKEN  = os.Getenv("BOT_TOKEN")
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	timestamp := int(time.Now().Unix())
	result := map[string]interface{}{}
	todoList := []string{"btcusd", "xrpusd", "ethusd", "ltcusd"}
	bot, _ := tgbotapi.NewBotAPI(BOT_TOKEN)
	msg := tgbotapi.NewMessageToChannel("-1001275593710", "Price Change:\n")
	for _, s := range todoList {
		target := fmt.Sprintf("%s/markets/bitfinex/%s/ohlc?periods=1800&before=%d&after=%d", ENDPOINT, s, timestamp, timestamp)

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
		info := rMap["result"]["1800"].([]interface{})[0]
		OHLC := map[string]float64{}
		fmt.Println(info)
		for index, value := range OHLCSTRUCT {
			OHLC[value] = info.([]interface{})[index].(float64)
		}

		// // Marshall that data into a map of AttributeValue object
		// av, err := dynamodbattribute.MarshalMap(OHLC)
		// if err != nil {
		// log.Println(err.Error())
		// }

		// sess := session.Must(session.NewSession())
		// svc := dynamodb.New(sess)
		// // Create DynamoDB client
		// if os.Getenv("AWS_SAM_LOCAL") == "true" {
		// log.Println("Serving in Local Environment...")
		// // svc = dynamodb.New(sess, aws.NewConfig().WithEndpoint("http://172.22.240.1:8000"))
		// svc = dynamodb.New(sess, aws.NewConfig().WithEndpoint("http://host.docker.internal:8000"))
		// }

		// input := &dynamodb.PutItemInput{
		// Item:      av,
		// TableName: aws.String("BTC_30m"),
		// }
		// _, err = svc.PutItem(input)

		// if err != nil {
		// log.Println("Got error calling PutItem:")
		// log.Println(err.Error())
		// return Response{StatusCode: 500, Body: err.Error()}, nil
		// }
		priceChange := (OHLC["closePrice"] - OHLC["openPrice"]) / OHLC["openPrice"]
		msg.Text += fmt.Sprintf("%s:\n Changed %.2f%%\n ClosePrice $%.2f\n", strings.ToUpper(s), priceChange*100, OHLC["closePrice"])
		result[s] = make(map[string]float64)
		for key, value := range OHLC {
			result[s].(map[string]float64)[key] = value
		}
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
