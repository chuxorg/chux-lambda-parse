package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	cfg "github.com/chuxorg/chux-lambda-parser/config"
	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
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
	cfg := cfg.LoadConfig()
	bucket := s3.New(
		s3.WithConfig(*cfg),
	)

	files, err := bucket.Download()
	if err != nil {
		panic(err)
	}

	parser := parsing.New(parsing.WithConfig(*cfg))
	for _, f := range files {
		parser.Parse(f)
	}
	return parseLambda.Parse(event.Input)
}

func main() {
	lambda.Start(parseHandler)
}
