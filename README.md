# firestore-cli

firestore-cli is a combined cli and terminal ui created to lookup firestore documents directly from your terminal.

For now you have to build from source yourself and moving it to a globally available folder.
```
go build -o fsctl cmd/*.go 

sudo chmod +x fsctl

sudo mv fsctl /usr/local/bin/fs

fsctl
```


