package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type IParseLambda interface {
	Parse(input string) (string, error)
}

type ParseLambda struct{}

var parseLambda IParseLambda = &ParseLambda{}

func (l *ParseLambda) Parse(input string) (string, error) {
	// Your parse implementation goes here
	return "", nil
}

type ParseEvent struct {
	Input string `json:"input"`
}

func parseHandler(ctx context.Context, event ParseEvent) (string, error) {
	return parseLambda.Parse(event.Input)
}

func main() {
	lambda.Start(parseHandler)
}
