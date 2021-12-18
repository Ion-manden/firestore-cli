# firestore-cli

firestore-cli is a combined cli and terminal ui created to lookup firestore documents directly from your terminal.

## Installation

## Homebrew
```
# add the repository
brew tab Ion-manden/firestore-cli

# install
brew install firestore-cli
```
### From source
```
go build -o fsctl cmd/*.go 

sudo chmod +x fsctl

sudo mv fsctl /usr/local/bin/fsctl

fsctl
```

### Using binary
Download the latest binary from the releases and place it in `/usr/local/bin/` or another exported PATH to allow usage in any folder.


## Usage
The cli has a couple different methods to help you work with firestore

##### get
`fsctl get <collection id> <document id>`
Will json data from the selected document and print it to stdout so you can pipe it along.
F.x. into jq to prettify it
`fsctl get tasks task1 | jq .`

##### ls
`fctl ls`
Simple list all collections in the GCP project.

##### count
`fctl count <collection id>`
Count all documents in a collection, due to the way firestore works this is a costly operation as you have to pay for document reads for each document in the collection.


##### tui
`fsctl tui`
Opens a terminal ui where you can explore your collections and documents.
