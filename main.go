package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/GeoffreyMetais/dlmanager/db"

	"github.com/ant0ine/go-json-rest/rest"
)

//Settings json object structure
var Settings struct {
	Root    string
	Port    string
	BaseURL string
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type directory struct {
	Path  string `json:"path"`
	Files []file `json:"files"`
}

type file struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDirectory"`
	Size  int64  `json:"size"`
}

type reqBody struct {
	Path string `json:"path"`
	User string `json:"user"`
}

func browseDir(w rest.ResponseWriter, req *rest.Request) {
	//basePath := "/mnt/hdd/usb/"
	//basePath := "/home/metais/Vidéos/"
	basePath := Settings.Root
	var request reqBody
	req.DecodeJsonPayload(&request)
	var path string
	if len(request.Path) > 0 {
		path = request.Path
	} else {
		path = basePath + req.PathParam("dir")
	}
	if !strings.HasPrefix(path, basePath) {
		rest.Error(w, "Permission denied", 503)
		return
	}

	fi, err := os.Stat(path)
	if err != nil {
		rest.NotFound(w, req)
		return
	}
	if fi.IsDir() {
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}
		filesInfo, _ := ioutil.ReadDir(path)
		//	files := [len(filesInfo)]string
		files := make([]file, len(filesInfo))
		for i := 0; i < len(filesInfo); i++ {
			files[i] = file{
				Name:  filesInfo[i].Name(),
				Path:  path + filesInfo[i].Name(),
				Size:  filesInfo[i].Size(),
				IsDir: filesInfo[i].IsDir(),
			}
		}
		dir := directory{
			Path:  path,
			Files: files,
		}
		/*w.Header().Set("Access-Control-Allow-Origin", "*")
		  w.Header().Set("Access-Control-Allow-Methods", "POST")*/
		w.WriteJson(&dir)
	} else {
		http.ServeFile(w.(http.ResponseWriter), req.Request, path)
	}
}

func download(w rest.ResponseWriter, req *rest.Request) {
	share := db.FindShare(req.PathParam("name"))
	fi, err := os.Stat(share.Path)
	if err != nil {
		rest.NotFound(w, req)
		return
	}
	w.Header().Add("Content-type", "application/octet-stream")
	//   w.Header().Add("Content-Type", "application/force-download")
	w.Header().Add("Content-Disposition", "attachment; filename="+fi.Name())
	w.Header().Add("Content-length", strconv.FormatInt(fi.Size(), 10))
	http.ServeFile(w.(http.ResponseWriter), req.Request, share.Path)
}

func listShares(w rest.ResponseWriter, req *rest.Request) {
	//     w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteJson(db.ListShares())
}

func add(w rest.ResponseWriter, req *rest.Request) {
	var newfile db.SharedFile
	req.DecodeJsonPayload(&newfile)
	if len(newfile.Path) > 0 && len(newfile.Name) > 0 {
		newfile.Link = Settings.BaseURL + "go/dl/" + newfile.Name
		db.Add(&newfile)
		w.WriteHeader(http.StatusOK)
	}
}

func remove(w rest.ResponseWriter, req *rest.Request) {
	filename, _ := url.QueryUnescape(req.PathParam("name"))
	db.Remove(filename)
	w.WriteHeader(http.StatusOK)
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

func main() {
	api := rest.NewApi()
	statusMw := &rest.StatusMiddleware{}
	api.Use(statusMw)
	api.Use(rest.DefaultDevStack...)
	defer db.PrepareDb().Close()
	router, err := rest.MakeRouter(
		rest.Get("/go/browse/#dir", browseDir),
		rest.Post("/go/browse", browseDir),
		rest.Get("/go/browse", browseDir),
		rest.Get("/go/dl/#name", download),
		rest.Post("/go/add", add),
		rest.Delete("/go/del/#name", remove),
		rest.Get("/go/list", listShares),
		rest.Get("/go/status", func(w rest.ResponseWriter, r *rest.Request) {
			w.WriteJson(statusMw.GetStatus())
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(Settings.Port, api.MakeHandler()))

}
