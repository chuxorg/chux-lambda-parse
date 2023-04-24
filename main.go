package main

import (
	"log"

	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
)

func main() {
	log.Println("Starting parse application")

	bucket := s3.New()

	log.Println("Downloading Files.")
	files, err := bucket.Download()
	if err != nil {
		log.Printf("Error downloading files: %v", err)
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
	return
}
