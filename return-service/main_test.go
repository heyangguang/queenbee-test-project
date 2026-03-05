package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturnHandler(t *testing.T) {
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
