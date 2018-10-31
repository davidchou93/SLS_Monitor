package main

import (
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
	ENDPOINT = os.Getenv("CRYPTOWATCH")
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	log.Println(ENDPOINT)
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
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(rBody),
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
