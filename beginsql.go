package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Expense struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
}

type ExpenseRequest struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

var db *sql.DB

// ADD EXPENSE
func addExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", 405)
		return
	}

	var req ExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "NADOGRY JSON", 400)
		return
	}

	var id int
	_, err = db.Exec(context.Background(), "INSERT INTO expenses(description, amount, date) VALUES($1,$2,$3)", req.Description, req.Amount, time.Now())

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	exp := Expense{
		ID:          id,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        time.Now().Format("2006-01-02"),
	}

	json.NewEncoder(w).Encode(exp)
}

// GET ALL
func getExpenses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(context.Background(), "SELECT id, description, amount, date FROM expenses")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	//defer rows.Close()

	var expenses []Expense

	for rows.Next() {
		var e Expense
		var date time.Time

		err = rows.Scan(&e.ID, &e.Description, &e.Amount, &date)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		e.Date = date.Format("2006-01-02")
		expenses = append(expenses, e)
	}

	json.NewEncoder(w).Encode(expenses)
}

// DELETE
func deleteExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid method", 405)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	_, err = db.Exec(context.Background(), "DELETE FROM expenses WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
}

// UPDATE
func updateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid method", 405)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	var req ExpenseRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	_, err = db.Exec(context.Background(), "UPDATE expenses SET description=$1, amount=$2 WHERE id=$3", req.Description, req.Amount, id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"updated": true})
}

// SUMMARY
func summary(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")

	w.Header().Set("Content-Type", "application/json")

	if monthStr == "" {
		var total float64
		db.QueryRow("SELECT COALESCE(SUM(amount),0) FROM expenses").Scan(&total)

		json.NewEncoder(w).Encode(map[string]float64{"total": total})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Invalid month", 400)
		return
	}

	var total float64
	db.QueryRow(context.Background(), `SELECT SUM(amount) FROM expenses WHERE EXTRACT(MONTH FROM date) = $1	`, month).Scan(&total)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"month": month,
		"total": total,
	})
}

// ROUTING
func expensesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getExpenses(w, r)
	case http.MethodPost:
		addExpense(w, r)
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func expenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteExpense(w, r)
	case http.MethodPut:
		updateExpense(w, r)
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func main() {
	var err error
	//db = connectDB("postgres://postgres:123545@localhost:5432/exam_db")
	//defer db.Close(context.Background())
	connStr := "user=postgres password=12345678 dbname=expenses_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("ISHLEDI")

	http.HandleFunc("/expenses", expensesHandler)
	http.HandleFunc("/expenses/", expenseByIDHandler)
	http.HandleFunc("/expenses/summary", summary)

	http.ListenAndServe(":8080", nil)
}

// func connectDB(config string) *pgx.Conn {
// 	conn, err = pgx.Connect(context.Background(), config)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 		os.Exit(1)
// 	}
// 	return conn
// }
