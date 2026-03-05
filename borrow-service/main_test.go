package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBorrowSuccess(t *testing.T) {
	req := BorrowRequest{UserID: 1, BookID: 1}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/api/borrow", bytes.NewReader(body))
	w := httptest.NewRecorder()

	borrowHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestBorrowOutOfStock(t *testing.T) {
	req := BorrowRequest{UserID: 1, BookID: 2}
	body, _ := json.Marshal(req)

	r := httptest.NewRequest(http.MethodPost, "/api/borrow", bytes.NewReader(body))
	w := httptest.NewRecorder()

	borrowHandler(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
