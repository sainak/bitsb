package handler

import (
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	t.Log("Test Ping")

	r := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	Ping(w, r)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}
