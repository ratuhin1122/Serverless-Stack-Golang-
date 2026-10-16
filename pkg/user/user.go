package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)
var (
	ErrorFailedToFatchRecord = "Error Failed to fatch Record"
	
	
)
type User struct{
	email   string `json:"email"`
	firstName string `json:"firstName"`
	lastName  string `json:"lastName"`
}


func FatchUser(email string, db dynamodbiface.DynamoDBAPI, tableName string) (*User, error){
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName)
	}
	result, err := db.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFatchRecord)
	}

	item := new(User)
	err =  dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToFatchRecord)
	}
	return item, nil

}

func FatchUsers(tablename string, db dynamodbiface.DynamoDBAPI) (*[]User,error){
	input := &dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}
	result,err :=db.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFatchRecord)
	}
	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, errors.New(ErrorFailedToFatchRecord)
	}
	return items, nil



}