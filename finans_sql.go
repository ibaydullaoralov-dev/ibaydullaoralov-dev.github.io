package main

import (
	"encoding/json"
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

var expenses []Expense
var nextID = 1

//  ADD EXPENSE
func addExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req ExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	exp := Expense{
		ID:          nextID,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        time.Now().Format("2006-01-02"),
	}

	nextID++
	expenses = append(expenses, exp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exp)
}

//  GET ALL
func getExpenses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

//  DELETE
func deleteExpense(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, e := range expenses {
		if e.ID == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

//  UPDATE
func updateExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req ExpenseRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for i, e := range expenses {
		if e.ID == id {
			expenses[i].Description = req.Description
			expenses[i].Amount = req.Amount

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(expenses[i])
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// 📊 SUMMARY
func summary(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")
	total := 0.0

	w.Header().Set("Content-Type", "application/json")

	if monthStr == "" {
		for _, e := range expenses {
			total += e.Amount
		}
		json.NewEncoder(w).Encode(map[string]float64{"total": total})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	for _, e := range expenses {
		t, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}

		if int(t.Month()) == month {
			total += e.Amount
		}
	}

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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func expenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteExpense(w, r)
	case http.MethodPut:
		updateExpense(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/expenses", expensesHandler)
	http.HandleFunc("/expenses/", expenseByIDHandler)
	http.HandleFunc("/expenses/summary", summary)

	http.ListenAndServe(":8080", nil)
}