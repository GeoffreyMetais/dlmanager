package db

import (
	"encoding/json"
	"fmt"
	"os"

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

//SharedFile to expose download link for files
type SharedFile struct {
	gorm.Model
	Name string `json:"name"`
	Path string `json:"path"`
	Link string `json:"link"`
}

func init() {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("opening config file", err.Error())
	} else {
		defer configFile.Close()
		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(&Settings); err != nil {
			fmt.Println("parsing config file", err.Error())
		}
	}
}

//Add share to database
func Add(file *SharedFile) {
	database.Create(file)
}

//Remove corresponding share from database
func Remove(filename string) {
	database.Delete(SharedFile{}, "Name = ?", filename)
}

//FindShare returns SharedFile corresponding to filename
func FindShare(filename string) SharedFile {
	share := SharedFile{}
	database.Find(&share, "name = ?", filename)
	return share
}

//PrepareDb setups the sqlite database
func PrepareDb() *gorm.DB {
	var err error
	database, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(&SharedFile{})
	return database
}

//ListShares lists all shares
func ListShares() []SharedFile {
	shares := []SharedFile{}
	database.Find(&shares)
	return shares
}
