package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BorrowRecord struct {
	ID         int        `json:"id"`
	BookID     int        `json:"book_id"`
	UserID     int        `json:"user_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"`
	IsOverdue  bool       `json:"is_overdue"`
}

var (
	borrowRecords = make(map[int]*BorrowRecord)
	bookStock     = make(map[int]int)
	mu            sync.RWMutex
)

func init() {
	borrowRecords[1] = &BorrowRecord{
		ID:         1,
		BookID:     1,
		UserID:     1,
		BorrowDate: time.Now().AddDate(0, 0, -35),
		DueDate:    time.Now().AddDate(0, 0, -5),
	}
	bookStock[1] = 5
}

func main() {
	http.HandleFunc("/api/borrows/", returnHandler)
	http.ListenAndServe(":8083", nil)
}

func returnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/borrows/")
	path = strings.TrimSuffix(path, "/return")
	id, err := strconv.Atoi(path)
	if err != nil {
		respondError(w, "Invalid borrow ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	record, exists := borrowRecords[id]
	if !exists {
		respondError(w, "Borrow record not found", http.StatusNotFound)
		return
	}

	if record.ReturnDate != nil {
		respondError(w, "Book already returned", http.StatusBadRequest)
		return
	}

	now := time.Now()
	record.ReturnDate = &now
	record.IsOverdue = now.After(record.DueDate)
	bookStock[record.BookID]++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
