package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	token string
	db    *sql.DB
)

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
	username := os.Getenv("username")
	password := os.Getenv("password")
	host := os.Getenv("host")
	db_port := os.Getenv("db_port")
	cnvDbPort, err := strconv.Atoi(db_port)
	if err != nil {
		fmt.Println("db_port isnt int")
		panic(err.Error())
	}
	databaseName := os.Getenv("database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, cnvDbPort, databaseName)
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	db = database
}

func main() {
	templates, err := template.ParseFiles("./task2/template/index.html", "./task2/template/upload.html")
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", ServeHTML(templates))
	http.HandleFunc("/upload", UploadHandler(templates))

	fmt.Println("Hosted in localhost:8000")
	defer db.Close()
	http.ListenAndServe(":8000", nil)
}

func ServeHTML(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
				"token": token,
			})
		} else {
			w.WriteHeader(404)
		}
	}
}

func UploadHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
