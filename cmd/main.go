package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ratuhin1122/serverless-stack/pkg/handlers"
)

var dynaClient dynamodbiface.DynamoDBAPI

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}
	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tableName = "LambdaInGo"

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		result, err := handlers.CreateUser(req, dynaClient, tableName)
		if result != nil {
			return *result, err
		}
		return events.APIGatewayProxyResponse{}, err
	case "GET":
		result, err := handlers.GetUser(req, dynaClient, tableName)
		if result != nil {
			return *result, err
		}
		return events.APIGatewayProxyResponse{}, err
	case "PUT":
		result, err := handlers.UpdateUser(req, dynaClient, tableName)
		if result != nil {
			return *result, err
		}
		return events.APIGatewayProxyResponse{}, err
	case "DELETE":
		result, err := handlers.DeleteUser(req, dynaClient, tableName)
		if result != nil {
			return *result, err
		}
		return events.APIGatewayProxyResponse{}, err
	default:
		result, err := handlers.UnhandleMethod()
		if result != nil {
			return *result, err
		}
		return events.APIGatewayProxyResponse{}, err
	}
}