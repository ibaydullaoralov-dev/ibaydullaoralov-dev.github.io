package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

// CRUD (Create, Read, Update, Delete)

type AuthorResponse struct {
	AuthorID  int    `json:"author_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthYear int    `json:"born_year"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type AuthorCreateRequest struct {
	// AuthorID  int    `json:"author_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthYear int    `json:"born_year"`
}

func authors(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), "select author_id, first_name, last_name, birth_year from authors")
	if err != nil {
		fmt.Fprintln(w, "YALNYSLYK: ", err.Error())
	}
	var list []AuthorResponse
	for rows.Next() {
		var res AuthorResponse
		err = rows.Scan(&res.AuthorID, &res.FirstName, &res.LastName, &res.BirthYear)
		if err != nil {
			fmt.Fprintln(w, "YALNYSLYK: ", err.Error())
		}
		list = append(list, res)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func createAuthor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(ErrorResponse{"dine POST method isleyar", "400"})
		return
	}
	var req AuthorCreateRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "400"})
		return
	}
	birthYearStr := strconv.Itoa(req.BirthYear)
	_, err = db.Exec(context.Background(), "insert into authors (first_name, last_name, birth_year) values ( '"+req.FirstName+"', '"+req.LastName+"', "+birthYearStr+");")
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{err.Error(), "400"})
		return
	}
	json.NewEncoder(w).Encode(true)
}

var db *pgx.Conn

func main() {
	db = connectDB("postgres://postgres:@localhost:5432/library_db")
	defer db.Close(context.Background())
	http.HandleFunc("/authors", authors)
	http.HandleFunc("/authors/create", createAuthor)
	http.ListenAndServe(":8000", nil)
}

func connectDB(config string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}