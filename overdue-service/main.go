package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// BorrowRecord 借阅记录
type BorrowRecord struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	BookID     int       `json:"book_id"`
	BorrowDate time.Time `json:"borrow_date"`
	DueDate    time.Time `json:"due_date"`
	ReturnDate *time.Time `json:"return_date"`
}

// OverdueRecord 逾期记录
type OverdueRecord struct {
	BorrowID   int    `json:"borrow_id"`
	UserID     int    `json:"user_id"`
	BookID     int    `json:"book_id"`
	DueDate    string `json:"due_date"`
	OverdueDays int   `json:"overdue_days"`
}

// 模拟借阅记录数据
var borrowRecords = []BorrowRecord{
	{ID: 1, UserID: 1, BookID: 1, BorrowDate: time.Now().AddDate(0, 0, -40), DueDate: time.Now().AddDate(0, 0, -10), ReturnDate: nil},
	{ID: 2, UserID: 2, BookID: 2, BorrowDate: time.Now().AddDate(0, 0, -35), DueDate: time.Now().AddDate(0, 0, -5), ReturnDate: nil},
	{ID: 3, UserID: 1, BookID: 3, BorrowDate: time.Now().AddDate(0, 0, -20), DueDate: time.Now().AddDate(0, 0, 10), ReturnDate: nil},
}

func main() {
	http.HandleFunc("/borrows/overdue", allOverdueHandler)
	http.HandleFunc("/users/", userOverdueHandler)
	http.ListenAndServe(":8084", nil)
}

// allOverdueHandler 查询所有逾期记录
func allOverdueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	overdueRecords := getOverdueRecords(0)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overdueRecords)
}

// userOverdueHandler 查询指定用户的逾期记录
func userOverdueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析用户ID: /users/:id/overdue
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "overdue" {
		respondError(w, "Invalid path", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	overdueRecords := getOverdueRecords(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overdueRecords)
}

// getOverdueRecords 获取逾期记录，userID为0时返回所有用户
func getOverdueRecords(userID int) []OverdueRecord {
	var result []OverdueRecord
	now := time.Now()

	for _, record := range borrowRecords {
		// 已归还的不算逾期
		if record.ReturnDate != nil {
			continue
		}

		// 未到期的不算逾期
		if now.Before(record.DueDate) {
			continue
		}

		// 如果指定了用户ID，只返回该用户的记录
		if userID > 0 && record.UserID != userID {
			continue
		}

		// 计算逾期天数
		overdueDays := int(now.Sub(record.DueDate).Hours() / 24)

		result = append(result, OverdueRecord{
			BorrowID:    record.ID,
			UserID:      record.UserID,
			BookID:      record.BookID,
			DueDate:     record.DueDate.Format("2006-01-02"),
			OverdueDays: overdueDays,
		})
	}

	return result
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
