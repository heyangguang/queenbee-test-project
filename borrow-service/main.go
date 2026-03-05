package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// 图书库存
type Book struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Stock int    `json:"stock"`
}

// 借阅记录
type BorrowRecord struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	BookID     int       `json:"book_id"`
	BorrowDate time.Time `json:"borrow_date"`
	DueDate    time.Time `json:"due_date"`
}

type BorrowRequest struct {
	UserID int `json:"user_id"`
	BookID int `json:"book_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	books   = make(map[int]*Book)
	records = make(map[int]*BorrowRecord)
	mu      sync.RWMutex
	nextID  = 1
)

func init() {
	// 初始化测试数据
	books[1] = &Book{ID: 1, Title: "Go编程", Stock: 5}
	books[2] = &Book{ID: 2, Title: "Python实战", Stock: 0}
}

func main() {
	http.HandleFunc("/api/borrow", borrowHandler)
	http.ListenAndServe(":8081", nil)
}

func borrowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BorrowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	book, exists := books[req.BookID]
	if !exists {
		respondError(w, "Book not found", http.StatusNotFound)
		return
	}

	if book.Stock <= 0 {
		respondError(w, "Book out of stock", http.StatusBadRequest)
		return
	}

	// 创建借阅记录
	record := &BorrowRecord{
		ID:         nextID,
		UserID:     req.UserID,
		BookID:     req.BookID,
		BorrowDate: time.Now(),
		DueDate:    time.Now().AddDate(0, 0, 30), // 30天后归还
	}
	records[nextID] = record
	nextID++

	// 库存-1
	book.Stock--

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
