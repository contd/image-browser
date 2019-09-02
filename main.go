package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	s "strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

// Directory struct is for directory objects
type Directory struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modified"`
	Path    string    `json:"path"`
}

// Exif picture data struct
type Exif struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// Picture struct is for picture objects
type Picture struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Type    string    `json:"type"`
	ModTime time.Time `json:"modified"`
	Path    string    `json:"path"`
	Width   int       `json:"width"`
	Height  int       `json:"height"`
	Exif    Exif      `json:"exif"`
}

// Files struct stores the current directory's contents in separate arrays of directores or pictures
type Files struct {
	Directories []Directory
	Pictures    []Picture
}

// RootPath is the starting path of directories or pictures
var RootPath string

// Directories is an array of Directory objects in current path
var Directories []Directory

// Pictures is an array of Picture objects in current path
var Pictures []Picture

// CollectionPath is the name of the link to the collection being served
var CollectionPath = "collections"

// PageLimit is the number of items per page
var PageLimit = 10

// OriginsAllowed is for CORS and should leave the localhost there
var OriginsAllowed = []string{"*"}

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %s", err)
	}
	RootPath = path.Join(root, CollectionPath)
	Directories, Pictures = readPath(RootPath)

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	// CORS restricted
	// Allows requests from any `https://kumpf.io` or `http://localhost:*` origin
	// wth GET, PUT, POST or DELETE method.
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: OriginsAllowed,
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	router.Static("/collections", "./collections")
	router.Static("/", "./public")
	router.GET("/api", listRoot)
	router.Logger.Fatal(router.Start(":6969"))
}

func listRoot(c echo.Context) error {
	pathParam := c.QueryParam("path")
	pageParam, _ := strconv.Atoi(c.QueryParam("page"))
	if pageParam == 0 {
		pageParam = 1
	}
	files := getFiles(pathParam, pageParam)
	return c.JSON(http.StatusOK, files)
}

func getFiles(pathParam string, pageParam int) Files {
	if pathParam != "" && pathParam != "/" {
		Directories, Pictures = readPath(path.Join(RootPath, pathParam))
	} else {
		Directories, Pictures = readPath(RootPath)
	}
	sort.SliceStable(Directories, func(i, j int) bool {
		return Directories[i].Name < Directories[j].Name
	})
	sort.SliceStable(Pictures, func(i, j int) bool {
		return Pictures[i].Name < Pictures[j].Name
	})
	// Limit (Paging)
	Pictures = limitPics(pageParam, PageLimit)

	files := Files{
		Directories: Directories,
		Pictures:    Pictures,
	}
	return files
}

func readPath(fullPath string) ([]Directory, []Picture) {
	dir, err := os.Open(fullPath)
	if err != nil {
		log.Fatalf("failed to open dir: %s", err)
	}
	defer dir.Close()

	relPath := s.Replace(fullPath, RootPath, "", 1)
	items, err := dir.Readdir(0)
	if err != nil {
		log.Fatalf("failed to read dir: %s", err)
	}

	var dirs []Directory
	var pics []Picture

	for _, item := range items {
		if item.IsDir() {
			if item.Name() != "." && item.Name() != ".." {
				directory := Directory{
					Name:    item.Name(),
					Size:    item.Size(),
					ModTime: item.ModTime(),
					Path:    fmt.Sprintf("http://localhost:6969/api?path=%s", path.Join(relPath, item.Name())),
				}
				dirs = append(dirs, directory)
			}
		} else {
			if checkExtension(item.Name()) {
				pic, err := os.Open(path.Join(fullPath, item.Name()))
				if err != nil {
					log.Fatalf("failed to open pic: %v", item.Name())
				}

				mimeType := getMimeType(item.Name())
				image, _, err := image.DecodeConfig(pic)
				if err != nil {
					log.Fatalf("%s: %v", item.Name(), err)
				}

				exifInfo, exifErr := getExif(path.Join(fullPath, item.Name()))
				if exifErr != nil {
					exifInfo = Exif{Lat: 0, Long: 0}
				}
				picture := Picture{
					Name:    item.Name(),
					Size:    item.Size(),
					Type:    mimeType,
					ModTime: item.ModTime(),
					Path:    fmt.Sprintf("http://localhost:6969/collections%s", path.Join(relPath, item.Name())),
					Width:   image.Width,
					Height:  image.Height,
					Exif:    exifInfo,
				}
				pics = append(pics, picture)
			}
		}
	}

	return dirs, pics
}

func limitPics(p int, lim int) []Picture {
	e := p * lim
	s := e - lim
	if e > len(Pictures) {
		e = len(Pictures)
	}
	var rangePics []Picture

	for i := s; i < e; i++ {
		log.Printf("I: %v - LEN: %v", i, len(Pictures))
		rangePics = append(rangePics, Pictures[i])
	}

	return rangePics
}

func randPics(n int) []Picture {
	var randPics []Picture
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		pick := Pictures[rand.Intn(len(Pictures))]
		randPics = append(randPics, pick)
	}

	return randPics
}

func widthOrHeight(w int, h int) string {
	if w > h {
		return "width"
	}
	return "height"
}

func checkExtension(fname string) bool {
	re := regexp.MustCompile(`.*\.(?:jpg|gif|png|tif)$`)
	if re.MatchString(fname) {
		return true
	}
	return false
}

func getMimeType(fname string) string {
	extn := filepath.Ext(fname)
	mtype := mime.TypeByExtension(extn)
	return mtype
}

func getExif(fname string) (Exif, error) {
	exif.RegisterParsers(mknote.All...)

	f, err := os.Open(fname)
	if err != nil {
		log.Printf("Failed to open for exif [%s]: %v", fname, err)
		return Exif{}, err
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Printf("Failed to decode exif [%s]: %v", fname, err)
		return Exif{}, err
	}

	lat, long, _ := x.LatLong()

	return Exif{
		Lat:  lat,
		Long: long,
	}, nil
}
