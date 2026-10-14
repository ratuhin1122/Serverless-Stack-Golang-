package main

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ratuhin1122/serverless-stack/pkg/handlers"
)

var ( 
	db dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}
	db := dynamodb.New(awsSession)
	lambda.Start(handler)
	
}

const tableName = "LambdaInGo"
func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		handlers.CreateUser(req, db, tableName)
	case "GET":
		handlers.GetUser(req, db, tableName)
	case "PUT":
		handlers.UpdateUser(req, db, tableName)
	case "DELETE":
		handlers.DeleteUser(req, db, tableName)
	default:
		return handlers.UnhandleMethod()
	}
}