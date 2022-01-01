/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count documents in a collection",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("count needs exactly 1 arguments, {collection id}")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// start firestore client
		ctx := context.Background()
		conf := &firebase.Config{ProjectID: projectId}
		sa := option.WithCredentialsFile(credentialsFile)
		fba, err := firebase.NewApp(ctx, conf, sa)
		if err != nil {
			log.Fatalln(err)
		}

		client, err := fba.Firestore(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		defer client.Close()

		// Ask for confirmation
		fmt.Println("WARNING - Costly operation")
		fmt.Println("As firestore does not have a count method, the only way to count documents is to read them all and count the results")
		fmt.Println("This means GCP will charge you for document reads for every document in your collection")
		fmt.Println("Are you sure you want to proceed? (type YES to continue)")
		var response string
		_, err = fmt.Scanln(&response)
		if err != nil {
			log.Fatal(err)
		}
		if response != "YES" {
			fmt.Println("Operation canceled")
			return
		}

		docs, err := client.Collection(args[0]).Select(firestore.DocumentID).Documents(context.Background()).GetAll()
		if err != nil {
			log.Println(err)
			return
		}
		count := len(docs)
		fmt.Println(fmt.Sprintf("Total of %d documents in %v", count, args[0]))
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}
