package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func connectDB(config string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect: %v\n", err)
		os.Exit(1)
	}
	return conn
}

func authors(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(),
		"select author_id, first_name, last_name, birth_year from authors")
	if err != nil {
		fmt.Fprintln(w, "YALNYSLYK:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var firstName, lastName string
		var id, birthYear int

		err := rows.Scan(&id, &firstName, &lastName, &birthYear)
		if err != nil {
			fmt.Fprintln(w, "YALNYSLYK:", err)
			return
		}

		fmt.Fprintln(w, "NETIJE:", id, firstName, lastName, birthYear)
	}
}

func hosh(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,
		"Good bye World",
		"Method:", r.Method,
		"Header:", r.Header,
		"Host:", r.Host,
		"Path:", r.URL.Path,
		"RemoteAddr:", r.RemoteAddr,
		"RequestURI:", r.RequestURI,
	)
}

func main() {
	db = connectDB("postgres://ibaydulla:0000@localhost:5432/library_db")
	defer db.Close(context.Background())

	http.HandleFunc("/authors", authors)
	http.HandleFunc("/books", authors)
	http.HandleFunc("/hosh", hosh)
	http.HandleFunc("/bye", hosh)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
