package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
)

type ParseLambda struct{}

func (l *ParseLambda) Parse(ctx context.Context, input string) (string, error) {
	
	log.Println("Entering Parse Lambda")

	bucket := s3.New()

	log.Println("Downloading Files.")	
	files, err := bucket.Download()
	if err != nil {
		return "", err
	}
	log.Println("Files downloaded.")

	log.Println("Parsing Files.")
	parser := parsing.New()
	for _, f := range files {
		
		parser.Parse(f)
	}
	log.Println("Files Parsed.")

	log.Println("Saving Files to Mongo.")
	filesInterface := make([]interface{}, len(files))
	for i, file := range files {
		filesInterface[i] = file
	}
	file := s3.File{}
	file.Save(filesInterface)
	log.Println("Files Saved")
	
	return "", nil

}

type ParseEvent struct {
	Input string `json:"input"`
}

func parseHandler(ctx context.Context, event ParseEvent) (string, error) {
	log.Println("Setting up Parse Lambda")
	pl := ParseLambda{}
	s, err := pl.Parse(ctx, event.Input)
	if err != nil {
		log.Printf("Error in parseHandler: %v", err)
		return "", err
	}
	return s, nil
}

func main() {
	
	log.Println("Starting parse lambda")
	lambda.Start(parseHandler)
}
