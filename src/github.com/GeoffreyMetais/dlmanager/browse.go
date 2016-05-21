package main

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"net/http"
        "encoding/json"
	"strings"
	"os"
        "strconv"
)

var settings struct {
    Root string
    Port string
    BaseUrl string
}

var filesCollection struct {
   Pool map[string]SharedFile 
}

type SharedFile struct {
    Name string
    Path string
    Link string
}

type User struct {
	Id   string
	Name string
}

type Directory struct {
	Path  string
	Files []File
}

type File struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
}

type ReqBody struct {
	Path string
	User string
}

func browseDir(w rest.ResponseWriter, req *rest.Request) {
	//basePath := "/mnt/hdd/usb/"
	//basePath := "/home/metais/Vidéos/"
        basePath := settings.Root
	var request ReqBody
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
            files := make([]File, len(filesInfo))
            for i := 0; i < len(filesInfo); i++ {
                    files[i] = File{
                            Name:  filesInfo[i].Name(),
                            Path:  path + filesInfo[i].Name(),
                            Size:  filesInfo[i].Size(),
                            IsDir: filesInfo[i].IsDir(),
                    }
            }
            dir := Directory{
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
    fmt.Println("Download ")
    var filename = req.PathParam("name")
    fmt.Println("Download ", filename)
    fmt.Println("path ", filesCollection.Pool[filename].Path)
    fi, err := os.Stat(filesCollection.Pool[filename].Path)
    if err != nil {
        rest.NotFound(w, req)
        return
    }
    fmt.Println("stats ", fi.Mode())
    fmt.Println("stats name ", fi.Name())
    fmt.Println("formatted size ", strconv.FormatInt(fi.Size(),10))
    w.Header().Add("Content-type", "application/octet-stream");
 //   w.Header().Add("Content-Type", "application/force-download")
    w.Header().Add("Content-Disposition", "attachment; filename="+fi.Name())
    w.Header().Add("Content-length", strconv.FormatInt(fi.Size(),10))
    http.ServeFile(w.(http.ResponseWriter), req.Request, filesCollection.Pool[filename].Path)
}

func listShares(w rest.ResponseWriter, req *rest.Request) {
//     w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteJson(&filesCollection.Pool)
}

func init(){
    configFile, err := os.Open("config.json")
    if err != nil {
        fmt.Println("opening config file", err.Error())
    } else {
        defer configFile.Close();
        jsonParser := json.NewDecoder(configFile)
        if err = jsonParser.Decode(&settings); err != nil {
            fmt.Println("parsing config file", err.Error())
        }
    }
    ReadCollection()
}

func test() {
//     ReadCollection()
    for filename := range filesCollection.Pool {
        fmt.Println("name", filename)
        fmt.Println("path", filesCollection.Pool[filename].Path)
        fmt.Println("Link", filesCollection.Pool[filename].Link)
    }
//     filesCollection.Pool["troisème"] = SharedFile{"DBZ","/home/metais/Vidéos/[DB-Z.com] Dragon Ball Z Battle of Gods [720p][VOSTFR].mp4"}
//     SaveCollection()
//     fmt.Println("Collection saved")
}

func main() {
	handler := rest.ResourceHandler{
            PreRoutingMiddlewares: []rest.Middleware{
                        &rest.CorsMiddleware{
                                RejectNonCorsRequests: false,
                                OriginValidator: func(origin string, request *rest.Request) bool {
                                        return true //origin == "http://localhost:8000"
                                },
                                AllowedMethods: []string{"GET", "POST", "PUT"},
                                AllowedHeaders: []string{
                                        "Accept", "Content-Type", "X-Custom-Header", "Origin"},
                                AccessControlAllowCredentials: true,
                                AccessControlMaxAge:           3600,
                        },
                },
        }
	handler.SetRoutes(
		&rest.Route{"GET",  "/go/browse/:dir", browseDir},
		&rest.Route{"POST", "/go/browse", browseDir},
                &rest.Route{"GET",  "/go/browse", browseDir},
                &rest.Route{"GET",  "/go/dl/:name", download},
                &rest.Route{"POST", "/go/add", add},
                &rest.Route{"DELETE",  "/go/del/:name", remove},
                &rest.Route{"GET",  "/go/list", listShares},
	)
        //http.HandleFunc("/dl", download)
	http.ListenAndServe(settings.Port, &handler)
}
