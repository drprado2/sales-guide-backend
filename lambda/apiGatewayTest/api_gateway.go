package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"net/http"
)

type MyEvent struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type MyEventOutput struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func HandleLambdaEvent(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var resp MyEvent
	err := json.Unmarshal([]byte(req.Body), &resp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	response := MyEventOutput{
		Name: "Olá " + resp.Name,
		Age:  fmt.Sprintf("Sua ideia é %v", resp.Age),
	}
	js, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
		Headers: map[string]string{
			"x-cid":        uuid.NewString(),
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
