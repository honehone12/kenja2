package response

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func BadRequest() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       "bad request",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
}

func InternalError() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "internal server error",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
}

func Json(s string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       s,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
