package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chuxorg/chux-parser/parsing"
	"github.com/chuxorg/chux-parser/s3"
)

type ParseLambda struct{}

func (l *ParseLambda) Parse(ctx context.Context, input string) (string, error) {

}

type ParseEvent struct {
	Input string `json:"input"`
}

func parseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Setting up Parse Lambda")
	pl := ParseLambda{}

	var event ParseEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := pl.Parse(r.Context(), event.Input)
	if err != nil {
		log.Printf("Error in parseHandler: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": s})
}

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

}
