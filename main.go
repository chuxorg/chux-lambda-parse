package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
)

type IParseLambda interface {
	Parse(input string) (string, error)
}

type ParseLambda struct{}

var parseLambda IParseLambda = &ParseLambda{}

func (l *ParseLambda) Parse(input string) (string, error) {
	bucket := s3.New()

	files, err := bucket.Download()
	if err != nil {
		panic(err)
	}

	parser := parsing.New()
	for _, f := range files {
		parser.Parse(f)
	}
	filesInterface := make([]interface{}, len(files))
	for i, file := range files {
		filesInterface[i] = file
	}

	file := s3.File{}
	file.Save(filesInterface)
	return "", nil
}

type ParseEvent struct {
	Input string `json:"input"`
}

func parseHandler(ctx context.Context, event ParseEvent) (string, error) {
	log.Default().Println("Setting up Parse Lambda")
	pl := ParseLambda{}
	s, err := pl.Parse(event.Input)
	if err != nil {
		return s, err
	}
	return parseLambda.Parse(event.Input)
}

func main() {
	taskRoot := os.Getenv("LAMBDA_TASK_ROOT")
	if taskRoot != "" {
		err := os.Chdir(taskRoot)
		if err != nil {
			log.Fatalf("Failed to change working directory: %v", err)
		}
	}
	log.Default().Println("Starting parse lambda")
	lambda.Start(parseHandler)
}
