package db

import (
	"net/url"

	"github.com/jinzhu/gorm"
	//preload sqlite driver
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

var database *gorm.DB

//Exposed shares
type SharedFile struct {
	gorm.Model
	Name string `json:"name"`
	Path string `json:"path"`
	Link string `json:"link"`
}

func Add(file *SharedFile) {
	database.Create(file)
}

func Remove(filename string) {
	database.Delete(SharedFile{}, "Name = ?", filename)
}

func FindShare(filename string) SharedFile {
	name, _ := url.QueryUnescape(filename)
	share := SharedFile{Name: name, Path: "", Link: ""}
	database.Find(&share)
	return share
}

func PrepareDb() *gorm.DB {
	var err error
	database, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
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
