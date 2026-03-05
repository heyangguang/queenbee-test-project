package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// OverdueRecord 逾期记录
type OverdueRecord struct {
	BorrowID   int       `json:"borrow_id"`
	BookTitle  string    `json:"book_title"`
	Username   string    `json:"username"`
	DueDate    time.Time `json:"due_date"`
	OverdueDays int      `json:"overdue_days"`
}

func main() {
	http.HandleFunc("/api/overdue", overdueHandler)
	http.ListenAndServe(":8082", nil)
}

// overdueHandler 查询逾期图书列表
func overdueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 模拟逾期数据（实际应从数据库查询）
	now := time.Now()
	records := []OverdueRecord{
		{
			BorrowID:   1,
			BookTitle:  "Go编程",
			Username:   "user1",
			DueDate:    now.AddDate(0, 0, -5),
			OverdueDays: 5,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
