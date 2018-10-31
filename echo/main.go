package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

var (
	ENDPOINT   = os.Getenv("CRYPTOWATCH")
	OHLCSTRUCT = [6]string{"CloseTime", "OpenPrice", "HighPrice", "LowPrice", "ClosePrice", "Volume"}
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	timestamp := int(time.Now().Unix())
	target := fmt.Sprintf("%s/markets/bitfinex/btcusd/ohlc?periods=1800&before=%d&after=%d", ENDPOINT, timestamp, timestamp)
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
	for index, value := range info.([]interface{}) {
		OHLC[OHLCSTRUCT[index]] = value.(float64)
	}
	result, _ := json.Marshal(OHLC)
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(result),
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
