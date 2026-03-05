package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOverdueHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/overdue", nil)
	w := httptest.NewRecorder()

	overdueHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 200, 得到 %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("期望 Content-Type application/json, 得到 %s", contentType)
	}
}
