package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// ReturnRequest 归还请求
type ReturnRequest struct {
	BorrowID int `json:"borrow_id"`
}

// ReturnResponse 归还响应
type ReturnResponse struct {
	Message string `json:"message"`
}

// BorrowRecord 借阅记录（模拟数据）
var borrowRecords = map[int]map[string]interface{}{
	1: {"book_id": 1, "user_id": 1, "borrow_date": "2026-02-01", "return_date": nil},
}

// BookStock 图书库存（模拟数据）
var bookStock = map[int]int{
	1: 5,
}

// 并发安全锁
var mu sync.RWMutex

func main() {
	http.HandleFunc("/api/return", returnHandler)
	http.ListenAndServe(":8083", nil)
}

// returnHandler 处理图书归还
func returnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReturnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 输入验证：borrow_id 必须大于 0
	if req.BorrowID <= 0 {
		respondError(w, "Invalid borrow_id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	record, exists := borrowRecords[req.BorrowID]
	if !exists {
		respondError(w, "Borrow record not found", http.StatusNotFound)
		return
	}

	if record["return_date"] != nil {
		respondError(w, "Book already returned", http.StatusBadRequest)
		return
	}

	// 更新归还日期
	record["return_date"] = time.Now().Format("2006-01-02")

	// 库存+1
	bookID := record["book_id"].(int)
	bookStock[bookID]++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ReturnResponse{Message: "Book returned successfully"})
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
