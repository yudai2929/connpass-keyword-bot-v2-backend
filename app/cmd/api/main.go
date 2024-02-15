package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	log.Printf("echo cold start")

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(BasicAuth())
	e.GET("/api/healthcheck", Healthcheck)

	echoLambda = echoadapter.New(e)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}

type HealthcheckMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Healthcheck(c echo.Context) error {
	msg := &HealthcheckMessage{
		Status:  http.StatusOK,
		Message: "Success to connect echo",
	}
	res, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return c.String(http.StatusOK, string(res))
}

var (
	id = "test"
	pw = "pass"
)

func BasicAuth() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username string, password string, context echo.Context) (bool, error) {
		if username == id && password == pw {
			return true, nil
		}
		return false, nil
	})
}
