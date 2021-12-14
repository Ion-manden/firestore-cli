package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	var projectId string
	flag.StringVar(&projectId, "project", os.Getenv("GCP_PROJECT_ID"), "project id of firestore")
	flag.Parse()

	credentialsFile := os.Getenv("GCP_CREDENTIALS_FILE")
	if credentialsFile == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln(err)
		}
		credentialsFile = fmt.Sprint(homeDir, "/.config/gcloud/application_default_credentials.json")
	}

	if projectId == "" {
		fmt.Println(projectId)
		fmt.Println("Missing project id, please set you GCP_PROJECT_ID env varibale or pass it as a flag {-project <project id>}")
		return
	}

	if credentialsFile == "" {
		fmt.Println(credentialsFile)
		fmt.Println("Missing path for GCP credentials file, please set you GCP_CREDENTIALS_FILE env varibale")
		return
	}

	cmdArgs := flag.Args()

	// if no arguments, show help
	if len(cmdArgs) == 0 {
		fmt.Println("Firestore terminal ui")
		fmt.Println("Use tui to start the terminal ui")
		fmt.Println("Use ls to list available collections")
		fmt.Println("Use get {collection id} {document id} to get content of a document")

		return
	}

	// start firestore client
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectId}
	sa := option.WithCredentialsFile(credentialsFile)
	fba, err := firebase.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = fba.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	action := cmdArgs[0]

	switch action {
	case "tui":
		startTui(client)

	case "ls":
		cs, err := client.Collections(context.Background()).GetAll()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("Collections")
		fmt.Println()
		for _, c := range cs {
			fmt.Println("id:", c.ID)
			fmt.Println("path:", c.Path)
			fmt.Println()
		}

	case "get":
		// if two arguments, get doc from collection
		if len(cmdArgs) != 3 {
			fmt.Println("get needs exactly 2 arguments, {collection id} and {document id}")
			return
		}

		doc, err := client.Collection(cmdArgs[1]).Doc(cmdArgs[2]).Get(context.Background())
		if err != nil {
			log.Println(err)
			return
		}

		j, err := json.Marshal(doc.Data())
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(j))
	default:
		fmt.Println("argument not recognized")
		fmt.Println("Use tui to start the terminal ui")
		fmt.Println("Use ls to list available collections")
		fmt.Println("Use get {collection id} {document id} to get content of a document")
	}
}
