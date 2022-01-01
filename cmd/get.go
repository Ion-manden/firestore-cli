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
	"encoding/json"
	"errors"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"g"},
	Short:   "Easy way to get data from documents in firestore",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("get needs exactly 2 arguments, {collection id} and {document id}")
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

		doc, err := client.Collection(args[0]).Doc(args[1]).Get(context.Background())
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
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
