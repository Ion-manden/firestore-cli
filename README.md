# firestore-cli

firestore-cli is a combined cli and terminal ui created to lookup firestore documents directly from your terminal.

## Installation

### From source
```
go build -o fsctl cmd/*.go 

sudo chmod +x fsctl

sudo mv fsctl /usr/local/bin/fsctl

fsctl
```

### Using binary
Download the latest binary from the releases and place it in `/usr/local/bin/` or another exported PATH to allow usage in any folder.

