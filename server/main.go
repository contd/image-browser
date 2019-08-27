package main

import (
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	s "strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Directory struct is for directory objects
type Directory struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modified"`
	Path    string    `json:"path"`
}

// Picture struct is for picture objects
type Picture struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modified"`
	Path    string    `json:"path"`
	Width   int       `json:"width"`
	Height  int       `json:"height"`
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

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %s", err)
	}
	RootPath = path.Join(root, "collections")
	Directories, Pictures = readPath(RootPath)

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"widthOrHeight": widthOrHeight,
	})
	router.LoadHTMLFiles("templates/index.html")
	router.Static("/public", "./public")
	router.Static("/collections", "./collections")
	router.StaticFile("/favicon.ico", "./public/favicon.ico")
	router.GET("/", showIndex)
	router.GET("/api", listRoot)
	router.Run(":6969")
}

func widthOrHeight(w int, h int) string {
	if w > h {
		return "width"
	} else {
		return "height"
	}
}

func showIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Photo Serve",
	})
}

func listRoot(c *gin.Context) {
	pathParam := c.DefaultQuery("path", "")
	pageParam, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	files := getFiles(pathParam, pageParam)
	c.JSON(http.StatusOK, files)
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
	Pictures = limitPics(pageParam, 5)

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
	// take fullPath and remove os.Getwd() to make relPath
	relPath := s.Replace(fullPath, RootPath, "", 1)

	items, _ := dir.Readdir(0)

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
				image, _, err := image.DecodeConfig(pic)
				if err != nil {
					log.Fatalf("%s: %v", item.Name(), err)
				}
				picture := Picture{
					Name:    item.Name(),
					Size:    item.Size(),
					ModTime: item.ModTime(),
					Path:    fmt.Sprintf("http://localhost:6969/collections%s", path.Join(relPath, item.Name())),
					Width:   image.Width,
					Height:  image.Height,
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

func checkExtension(fname string) bool {
	if path.Ext(fname) == ".gif" || path.Ext(fname) == ".GIF" {
		return true
	} else {
		return false
	}
}
