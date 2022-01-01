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
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:     "ui",
	Aliases: []string{"tui"},
	Short:   "Open terminal ui to browse your collections",
	Long:    `The terminal ui is created to quickly browse your collections and documents.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		// Start the application.
		app = tview.NewApplication()
		finder(client)
		if err := app.Run(); err != nil {
			fmt.Printf("Error running application: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
}

type currenctSelectedCollection struct {
	i            int
	collectionId string
	t            string
	r            rune
}

var (
	client             *firestore.Client
	app                *tview.Application // The tview application.
	collections        *tview.List
	documents          *tview.List
	data               *tview.TextView
	input              *tview.InputField
	currenctCollection currenctSelectedCollection
)

func finder(client *firestore.Client) {

	// Create the basic objects.
	collections = tview.NewList().ShowSecondaryText(false)
	collections.SetBorder(true).SetTitle("Collections")
	documents = tview.NewList()
	documents.ShowSecondaryText(false)
	documents.SetBorder(true).SetTitle("Documents")
	data = tview.NewTextView().SetDynamicColors(true)
	data.SetBorder(true).SetTitle("Data")
	input = tview.NewInputField()
	input.SetBorder(true).SetTitle("Document lookup (/)")

	collections.SetDoneFunc(func() {
		app.Stop()
	})
	collections.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		}
		return event
	})

	documents.SetDoneFunc(func() {
		data.Clear()
		documents.Clear()
		app.SetFocus(collections)
	})
	documents.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		case '/':
			app.SetFocus(input)
			return tcell.NewEventKey(tcell.KeyRune, '/', tcell.ModNone)
		}
		return event
	})

	data.SetDoneFunc(func(key tcell.Key) {
		app.SetFocus(documents)
	})
	data.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case '/':
			app.SetFocus(input)
			return tcell.NewEventKey(tcell.KeyRune, '/', tcell.ModNone)
		}
		return event
	})

	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			app.SetFocus(documents)
		case tcell.KeyEnter:
			getDocument(input.GetText())
			app.SetFocus(data)
		}
	})

	// Create the layout.
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(collections, 0, 3, true).
			AddItem(documents, 0, 5, false).
			AddItem(input, 0, 1, false), 0, 1, true).
		AddItem(data, 0, 3, false)

	cs, err := client.Collections(context.Background()).GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	for _, c := range cs {
		c := c
		collections.AddItem(c.ID, c.Path, 0, nil)
	}
	collections.SetSelectedFunc(setDocumentList)
	// When the user selects a table, show its content.
	documents.SetSelectedFunc(func(i int, documentId string, t string, s rune) {
		// content(db, dbName, documentId)
		app.SetFocus(data)
	})

	app.SetRoot(flex, true)
}

func getDocument(documentId string) {
	// A collection was selected. Show all of its documents.
	data.Clear()
	c := client.Collection(currenctCollection.collectionId)

	d, err := c.Doc(documentId).Get(context.Background())
	if err != nil {
		data.SetText(fmt.Sprint("Error trying to get document", err))
		return
	}

	// data.SetText(fmt.Sprint(doc.Data()))
	j, err := json.MarshalIndent(d.Data(), "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	data.SetText(string(j))
}

func setDocumentList(i int, collectionId string, t string, r rune) {
	currenctCollection = currenctSelectedCollection{i, collectionId, t, r}

	// A collection was selected. Show all of its documents.
	data.Clear()
	documents.Clear()
	c := client.Collection(collectionId)
	ds := c.Documents(context.Background())

	for i := 1; i < 15; i++ {
		d, err := ds.Next()
		if err != nil {
			break
			// log.Fatalln("Error getting documents", err)
		}

		documents.AddItem(d.Ref.ID, collectionId, 0, nil)
	}
	ds.Stop()

	app.SetFocus(documents)

	// When the user navigates to a table, show its columns.
	documents.SetChangedFunc(func(i int, documentId string, collectionId string, s rune) {
		// A document was selected. Show its data.
		data.Clear()
		doc, err := client.Collection(collectionId).Doc(documentId).Get(context.Background())
		if err != nil {
			log.Fatalln("Error getting document on changed", err)
		}

		// data.SetText(fmt.Sprint(doc.Data()))
		j, err := json.MarshalIndent(doc.Data(), "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		data.SetText(string(j))
	})

	documents.SetCurrentItem(0) // Trigger the initial selection.

}
