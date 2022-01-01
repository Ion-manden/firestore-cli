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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
)

var cfgFile string
var projectId string
var credentialsFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "firestore-cli",
	Short: "firestore-cli is a combined cli and terminal ui created to lookup firestore documents directly from your terminal.",
	Long: `firestore-cli is created to help you work with GCP firestore directly from your terminal

  The cli contains methods for getting document data, listing collections and getting a collections document count.

  There is event a small terminal ui for browsing your collections.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Could not find home dir")
		return
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.firestore-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&projectId, "project", "p", os.Getenv("GCP_PROJECT_ID"), "Project id for your GCP project")
	rootCmd.MarkFlagRequired("project")
	rootCmd.PersistentFlags().StringVar(
		&credentialsFile,
		"credentials-file",
		fmt.Sprint(home, "/.config/gcloud/application_default_credentials.json"),
		"Path for you GCP credentials",
	)
	rootCmd.MarkFlagRequired("credentialsFile")
}
