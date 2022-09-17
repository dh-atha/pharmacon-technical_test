package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/joho/godotenv"
)

var token string
var db *sql.DB

type Metadata struct {
	ID          string
	Filename    string
	ContentType string
	Size        int
	Path        string
	CreatedAt   time.Time
}

func init() {
	if err := godotenv.Load("./task2/secret.env"); err != nil {
		panic(err.Error())
	}
	token = os.Getenv("token")

	database, err := sql.Open("mysql", "root:@tcp(localhost:3306)/pharmacon")
	if err != nil {
		panic(err.Error())
	}
	db = database
}

func main() {
	defer db.Close()
	templates, err := template.ParseFiles("./task2/template/index.html", "./task2/template/upload.html")
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
				"token": token,
			})
		} else {
			w.WriteHeader(404)
		}
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			auth := r.FormValue("auth")
			file, fileHeader, err := r.FormFile("data")
			maxSize := int64(8000000)
			if err != nil || auth != token || !strings.Contains(fileHeader.Header.Get("Content-Type"), "image") || fileHeader.Size > maxSize {
				w.WriteHeader(403)
				return
			}
			defer file.Close()

			uniqueID := uuid.NewString()
			destination := fmt.Sprint("./task2/images/", uniqueID, "-", fileHeader.Filename)
			destinationFile, err := os.Create(destination)
			if err != nil {
				fmt.Println(err.Error())
				w.WriteHeader(500)
				return
			}
			defer destinationFile.Close()
			io.Copy(destinationFile, file)

			metadata := Metadata{
				ID:          uniqueID,
				Filename:    fileHeader.Filename,
				ContentType: fileHeader.Header.Get("Content-Type"),
				Size:        int(fileHeader.Size),
				Path:        destination,
			}

			_, err = db.Exec("INSERT INTO metadata (id, filename, content_type, size, path) VALUES (?,?,?,?,?)", metadata.ID, metadata.Filename, metadata.ContentType, metadata.Size, metadata.Path)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}

			templates.ExecuteTemplate(w, "upload.html", map[string]interface{}{
				"filename": fileHeader.Filename,
			})
		} else {
			w.WriteHeader(404)
		}
	})

	fmt.Println("Hosted in localhost:8000")
	http.ListenAndServe(":8000", nil)
}
