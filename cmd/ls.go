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
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all firestore collections in your GCP project",
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
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
