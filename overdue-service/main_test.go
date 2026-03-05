package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllOverdueHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/borrows/overdue", nil)
	w := httptest.NewRecorder()

	allOverdueHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var records []OverdueRecord
	json.NewDecoder(w.Body).Decode(&records)

	if len(records) != 2 {
		t.Errorf("Expected 2 overdue records, got %d", len(records))
	}
}

func TestUserOverdueHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/1/overdue", nil)
	w := httptest.NewRecorder()

	userOverdueHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var records []OverdueRecord
	json.NewDecoder(w.Body).Decode(&records)

	if len(records) != 1 {
		t.Errorf("Expected 1 overdue record for user 1, got %d", len(records))
	}

	if len(records) > 0 && records[0].UserID != 1 {
		t.Errorf("Expected user_id 1, got %d", records[0].UserID)
	}
}
