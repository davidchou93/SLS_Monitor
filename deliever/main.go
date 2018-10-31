package main

import (
	"context"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	update := tgbotapi.Update{}
	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return Response{StatusCode: 400}, nil
	}

	// Marshall that data into a map of AttributeValue object
	av, err := dynamodbattribute.MarshalMap(update)

	// Create DynamoDB client
	sess, err := session.NewSession(&aws.Config{})
	svc := dynamodb.New(sess)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("BotUpdates"),
	}

	_, err = svc.PutItem(input)

	if err != nil {
		log.Println("Got error calling PutItem:")
		log.Println(err.Error())
		return Response{StatusCode: 500, Body: err.Error()}, nil
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
