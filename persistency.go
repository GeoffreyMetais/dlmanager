package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
)

func add(w rest.ResponseWriter, req *rest.Request) {
	var newfile sharedFile
	req.DecodeJsonPayload(&newfile)
	if len(newfile.Path) > 0 && len(newfile.Name) > 0 {
		fmt.Println("adding ", newfile.Name)
		newfile.Link = settings.BaseURL + "go/dl/" + newfile.Name
		filesCollection.Pool[newfile.Name] = newfile
		saveCollection()
		//         w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}
}

func remove(w rest.ResponseWriter, req *rest.Request) {
	var filename = req.PathParam("name")
	fmt.Println("removing ", filename)
	delete(filesCollection.Pool, filename)
	saveCollection()
	w.WriteHeader(http.StatusOK)
}

func readCollection() {
	collectionFile, err := os.Open("collection.json")
	if err != nil {
		fmt.Println("opening collection file", err.Error())
	} else {
		defer collectionFile.Close()
		var data = &filesCollection.Pool
		jsonParser := json.NewDecoder(collectionFile)
		if err = jsonParser.Decode(&data); err != nil {
			fmt.Println("parsing collection file", err.Error())
		}
	}
}

func saveCollection() {
	collectionFile, err := os.Create("collection.json")
	var data = &filesCollection.Pool
	if err != nil {
		fmt.Println("opening collection file", err.Error())
	} else {
		defer collectionFile.Close()
		fmt.Println("collection file opened, writing data")
		enc := json.NewEncoder(collectionFile)
		if err = enc.Encode(data); err != nil {
			fmt.Println("parsing collection file", err.Error())
		}
	}
}
