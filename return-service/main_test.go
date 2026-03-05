package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// resetTestData 重置测试数据，确保测试隔离
func resetTestData() {
	borrowRecords = map[int]map[string]interface{}{
		1: {"book_id": 1, "user_id": 1, "borrow_date": "2026-02-01", "return_date": nil},
	}
	bookStock = map[int]int{
		1: 5,
	}
}

func TestReturnHandler(t *testing.T) {
	resetTestData()
	// 测试成功归还
	req := ReturnRequest{BorrowID: 1}
	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// 测试重复归还
	r2 := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w2 := httptest.NewRecorder()

	returnHandler(w2, r2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for duplicate return, got %d", w2.Code)
	}
}

func TestReturnHandler_NotFound(t *testing.T) {
	resetTestData()
	req := ReturnRequest{BorrowID: 999}
	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestReturnHandler_InvalidJSON(t *testing.T) {
	resetTestData()
	r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestReturnHandler_InvalidMethod(t *testing.T) {
	resetTestData()
	r := httptest.NewRequest(http.MethodGet, "/api/return", nil)
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestReturnHandler_InvalidBorrowID(t *testing.T) {
	resetTestData()
	testCases := []int{0, -1, -100}

	for _, borrowID := range testCases {
		req := ReturnRequest{BorrowID: borrowID}
		body, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
		w := httptest.NewRecorder()

		returnHandler(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for borrow_id=%d, got %d", borrowID, w.Code)
		}
	}
}
