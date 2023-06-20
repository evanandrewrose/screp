package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/icza/screp/repparser"
)

type Event struct {
	Path string `json:"path"`
}

func GetFileFromURI(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	fmt.Printf("Handling request for %s!", event.Path)

	fileBytes, err := GetFileFromURI(event.Path)
	
	if err != nil {
		fmt.Printf("Error retrieving file from URI %s: %v\n", event.Path, err)
		return "could not get file", err
	}

	config, err := repparser.ParseConfig(fileBytes, repparser.Config{
		Commands: true,
	})

	if err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		return "could not parse config", err
	}

	config.Compute()
	
	selectedFields := map[string]interface{}{
		"Computed": config.Computed,
		"Header": config.Header,
	}

	jsonBytes, err := json.Marshal(selectedFields)
	if err != nil {
		fmt.Printf("Error serializing map to JSON: %v\n", err)
		return "could not serialize map to JSON", err
	}
	
	return string(jsonBytes), nil
}


func main() {
	lambda.Start(HandleRequest)
}