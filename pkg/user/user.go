package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ratuhin1122/serverless-stack/pkg/validators"
)
var (
	ErrorFailedToFatchRecord = "Error Failed to fatch Record"
	ErrorInvalidEmail        = "Error Invalid Email"
	ErrorFailedToUnmarshal   = "Error Failed to Unmarshal"
	ErrorFailedToMarshal     = "Error Failed to Marshal"
	ErrorFailedToCreateUser  = "Error Failed to Create User"
	ErrorInvalidUser         = "Error Invalid User"
	ErrorFailedToDeleteUser  = "Error Failed to Delete User"
	ErrorFailedToUpdateUser  = "Error Failed to Update User"
	ErrorUserAlreadyExist    = "Error User Already Exist"
	ErrorUserDoesNotExist    = "Error User Does Not Exist"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}


func FatchUser(email string, db dynamodbiface.DynamoDBAPI, tableName string) (*User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}
	result, err := db.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFatchRecord)
	}

	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
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
func CreateUser(req events.APIGatewayProxyRequest, db dynamodbiface.DynamoDBAPI, tableName string) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorFailedToUnmarshal)
	}
	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	// যদি user আগে থেকে থাকে তাহলে error দাও
	currentUser, _ := FatchUser(u.Email, db, tableName)
	if currentUser != nil && len(currentUser.Email) > 0 {
		return nil, errors.New(ErrorUserAlreadyExist)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorFailedToMarshal)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}
	_, err = db.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToCreateUser)
	}
	return &u, nil
}
func UpdateUser(req events.APIGatewayProxyRequest, db dynamodbiface.DynamoDBAPI, tableName string) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorFailedToUnmarshal)
	}

	// ✅ সঠিক লজিক: user না থাকলে update করা যাবে না
	currentUser, _ := FatchUser(u.Email, db, tableName)
	if currentUser == nil || len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	// firstName এবং lastName উভয়ই আপডেট করো
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(u.Email)},
		},
		UpdateExpression: aws.String("set firstName = :fn, lastName = :ln"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":fn": {S: aws.String(u.FirstName)},
			":ln": {S: aws.String(u.LastName)},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
	}
	_, err := db.UpdateItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToUpdateUser)
	}
	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, db dynamodbiface.DynamoDBAPI, tableName string) (*User, error) {
	email := req.QueryStringParameters["email"]
	if email == "" {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	// ✅ সঠিক লজিক: user না থাকলে delete করা যাবে না
	currentUser, _ := FatchUser(email, db, tableName)
	if currentUser == nil || len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {S: aws.String(email)},
		},
	}
	_, err := db.DeleteItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToDeleteUser)
	}
	return currentUser, nil
}