package db

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
)

var Settings struct {
	Root    string
	Port    string
	BaseURL string
}

var FilesCollection struct {
	Pool map[string]SharedFile `json:"sharesList"`
}

type SharedFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Link string `json:"link"`
}

func Add(w rest.ResponseWriter, req *rest.Request) {
	var newfile SharedFile
	req.DecodeJsonPayload(&newfile)
	if len(newfile.Path) > 0 && len(newfile.Name) > 0 {
		fmt.Println("adding ", newfile.Name)
		newfile.Link = Settings.BaseURL + "go/dl/" + newfile.Name
		FilesCollection.Pool[newfile.Name] = newfile
		saveCollection()
		//         w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	}
}

func Remove(w rest.ResponseWriter, req *rest.Request) {
	var filename = req.PathParam("name")
	fmt.Println("removing ", filename)
	delete(FilesCollection.Pool, filename)
	saveCollection()
	w.WriteHeader(http.StatusOK)
}

func ReadCollection() {
	collectionFile, err := os.Open("collection.json")
	if err != nil {
		fmt.Println("opening collection file", err.Error())
	} else {
		defer collectionFile.Close()
		var data = &FilesCollection.Pool
		jsonParser := json.NewDecoder(collectionFile)
		if err = jsonParser.Decode(&data); err != nil {
			fmt.Println("parsing collection file", err.Error())
		}
	}
}

func saveCollection() {
	collectionFile, err := os.Create("collection.json")
	var data = &FilesCollection.Pool
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
