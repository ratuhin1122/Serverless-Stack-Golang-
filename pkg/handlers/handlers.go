package handlers
import (
	"github.com/serverless-stack/user"
	"net/http"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

)
var ErrorMethodNotAllowed = " MEthod Not "
type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}


func GetUser(req events.APIGatewayProxyRequest, db dynamodbiface.DynamoDBAPI, tableName string) (*events.APIGatewayProxyResponse, error){

	email := req.QueryStringParameters["email"]
	if (len(email) > 0){
		result, err := user.FatchUser(email, db, tableName)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	
	}

	result, err := user.FetchUsers(db, tableName)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)

}
func CreateUser(req events.APIGatewayProxyRequest, db dynamodbiface.DynamoDBAPI, tableName string) (*events.APIGatewayProxyResponse, error){

}
func UpdateUser(req events.APIGatewayProxyRequest, db dynamodbiface.DynamoDBAPI, tableName string) (*events.APIGatewayProxyResponse, error){

}
func DeleteUser(req events.APIGatewayProxyRequest, db dynamodbiface.DynamoDBAPI, tableName string) (*events.APIGatewayProxyResponse, error){

}
func UnhandleMethod()(*events.APIGatewayProxyResponse, error){
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
	
}