package response

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func BadRequest() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusBadRequest,
		Body:       "bad request",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
}

func InternalError() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "internal server error",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}
}

func Json(s string) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       s,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "max-age=86400",
		},
	}
}
