package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReturnBook(t *testing.T) {
	mu.Lock()
	borrowRecords = make(map[int]*BorrowRecord)
	bookStock = make(map[int]int)
	borrowRecords[1] = &BorrowRecord{
		ID:         1,
		BookID:     1,
		UserID:     1,
		BorrowDate: time.Now().AddDate(0, 0, -10),
		DueDate:    time.Now().AddDate(0, 0, 20),
	}
	bookStock[1] = 5
	mu.Unlock()

	req := httptest.NewRequest(http.MethodPut, "/api/borrows/1/return", nil)
	w := httptest.NewRecorder()
	returnHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var record BorrowRecord
	json.NewDecoder(w.Body).Decode(&record)

	if record.ReturnDate == nil {
		t.Error("ReturnDate should not be nil")
	}

	if record.IsOverdue {
		t.Error("Should not be overdue")
	}

	if bookStock[1] != 6 {
		t.Errorf("Expected stock 6, got %d", bookStock[1])
	}
}

func TestReturnOverdueBook(t *testing.T) {
	mu.Lock()
	borrowRecords = make(map[int]*BorrowRecord)
	borrowRecords[2] = &BorrowRecord{
		ID:         2,
		BookID:     1,
		UserID:     1,
		BorrowDate: time.Now().AddDate(0, 0, -40),
		DueDate:    time.Now().AddDate(0, 0, -10),
	}
	bookStock[1] = 5
	mu.Unlock()

	req := httptest.NewRequest(http.MethodPut, "/api/borrows/2/return", nil)
	w := httptest.NewRecorder()
	returnHandler(w, req)

	var record BorrowRecord
	json.NewDecoder(w.Body).Decode(&record)

	if !record.IsOverdue {
		t.Error("Should be overdue")
	}
}

func TestReturnAlreadyReturned(t *testing.T) {
	now := time.Now()
	mu.Lock()
	borrowRecords = make(map[int]*BorrowRecord)
	borrowRecords[3] = &BorrowRecord{
		ID:         3,
		BookID:     1,
		UserID:     1,
		BorrowDate: time.Now().AddDate(0, 0, -10),
		DueDate:    time.Now().AddDate(0, 0, 20),
		ReturnDate: &now,
	}
	mu.Unlock()

	req := httptest.NewRequest(http.MethodPut, "/api/borrows/3/return", nil)
	w := httptest.NewRecorder()
	returnHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestReturnNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/borrows/999/return", nil)
	w := httptest.NewRecorder()
	returnHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
