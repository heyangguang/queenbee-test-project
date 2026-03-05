package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setup 重置测试数据
func setup() {
	borrowRecords = map[int]map[string]interface{}{
		1: {"book_id": 1, "user_id": 1, "borrow_date": "2026-02-01", "return_date": nil},
	}
	bookStock = map[int]int{
		1: 5,
	}
}

func TestReturnHandler(t *testing.T) {
	setup()

	req := ReturnRequest{BorrowID: 1}
	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestReturnHandlerDuplicate(t *testing.T) {
	setup()

	req := ReturnRequest{BorrowID: 1}
	body, _ := json.Marshal(req)

	// 第一次归还
	r1 := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w1 := httptest.NewRecorder()
	returnHandler(w1, r1)

	// 第二次归还（重复）
	r2 := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w2 := httptest.NewRecorder()
	returnHandler(w2, r2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for duplicate return, got %d", w2.Code)
	}
}

func TestReturnHandlerNotFound(t *testing.T) {
	setup()

	req := ReturnRequest{BorrowID: 999}
	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestReturnHandlerInvalidBody(t *testing.T) {
	setup()

	r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestReturnHandlerWrongMethod(t *testing.T) {
	setup()

	r := httptest.NewRequest(http.MethodGet, "/api/return", nil)
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestReturnHandlerNegativeID(t *testing.T) {
	setup()

	req := ReturnRequest{BorrowID: -1}
	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/api/return", bytes.NewReader(body))
	w := httptest.NewRecorder()

	returnHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for negative ID, got %d", w.Code)
	}
}
