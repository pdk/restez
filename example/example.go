package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pdk/dbscan"
	"github.com/pdk/restez"
)

func main() {

	doServer := flag.Bool("server", false, "run an HTTP server")
	doMigration := flag.Bool("migrate", false, "execute database migration (ie set up the db)")
	flag.Parse()

	if *doMigration {
		migrateDatabase()
	}

	if *doServer {
		runHTTPServer()
	}
}

func runHTTPServer() {
	http.HandleFunc("/new", restez.HandlePOST(addNewPerson))
	http.HandleFunc("/list", restez.HandleGET(listPeople))

	log.Printf("listening for HTTP requests on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func openDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	return db
}

func migrateDatabase() {

	db := openDatabase()

	_, err := db.Exec(`
		create table if not exists people (
			id integer primary key autoincrement,
			name varchar,
			age integer,
			favorite_food varchar
		)
	`)
	if err != nil {
		log.Fatalf("failed to create table people: %v", err)
	}
}

func listPeople(params map[string]string) ([]Person, error) {
	log.Printf("handling /list request")

	db := openDatabase()
	defer db.Close()

	rows, err := db.Query(`select name, age, favorite_food from people`)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query on people table: %v", err)
	}
	defer rows.Close()

	return dbscan.All[Person](rows)
}

func addNewPerson(req ExampleRequest) (ExampleResponse, error) {
	log.Printf("handling /new request")

	if req.Name == "" || req.Age == 0 || req.FavoriteFood == "" {
		return ExampleResponse{}, fmt.Errorf("the request was invalid. one or more values missing: %#v", req)
	}

	db := openDatabase()
	defer db.Close()
	_, err := db.Exec(`insert into people (name, age, favorite_food) values (?, ?, ?)`, req.Name, req.Age, req.FavoriteFood)
	if err != nil {
		return ExampleResponse{}, fmt.Errorf("failed to insert row into people: %v", err)
	}

	return ExampleResponse{
		Message:   "your post was received",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}

type Person struct {
	Name         string `json:"name"`
	Age          int    `json:"age"`
	FavoriteFood string `json:"favoriteFood"`
}

type ExampleRequest struct {
	Name         string `json:"name"`
	Age          int    `json:"age"`
	FavoriteFood string `json:"favoriteFood"`
}

type ExampleResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
