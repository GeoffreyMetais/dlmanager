package db

import (
	"net/http"
	"net/url"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	//preload sqlite driver
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

//Settings json object structure
var Settings struct {
	Root    string
	Port    string
	BaseURL string
}

var database *gorm.DB

//Exposed shares
type SharedFile struct {
	gorm.Model
	Name string `json:"name"`
	Path string `json:"path"`
	Link string `json:"link"`
}

func Add(w rest.ResponseWriter, req *rest.Request) {
	var newfile SharedFile
	req.DecodeJsonPayload(&newfile)
	if len(newfile.Path) > 0 && len(newfile.Name) > 0 {
		newfile.Link = Settings.BaseURL + "go/dl/" + newfile.Name
		database.Create(&newfile)
		w.WriteHeader(http.StatusOK)
	}
}

func Remove(w rest.ResponseWriter, req *rest.Request) {
	filename, _ := url.QueryUnescape(req.PathParam("name"))
	database.Delete(SharedFile{}, "Name = ?", filename)
	w.WriteHeader(http.StatusOK)
}

func FindShare(filename string) SharedFile {
	name, _ := url.QueryUnescape(filename)
	share := SharedFile{Name: name, Path: "", Link: ""}
	database.Find(&share)
	return share
}

func ReadCollection() *gorm.DB {
	var err error
	database, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	database.AutoMigrate(&SharedFile{})
	return database
}

//List all shares
func ListShares() []SharedFile {
	shares := []SharedFile{}
	database.Find(&shares)
	return shares
}
