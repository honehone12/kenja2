package lambda

import (
	"kenja2/endec"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func ResponseBadRequest() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       "bad request",
		Headers: map[string]string{
			"Content-Type": "text/plain; charset=utf8",
		},
	}
}

func ResponseInternalError() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "internal server error",
		Headers: map[string]string{
			"Content-Type": "text/plain; charset=utf8",
		},
	}
}

func ResponseOk[E endec.Encoder](e E, b []byte) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       e.String(b),
		Headers: map[string]string{
			"Content-Type": e.ContentType(),
		},
	}
}
