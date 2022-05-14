package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mtlynch/picoshare/v2/handlers"
	"github.com/mtlynch/picoshare/v2/store"
	"github.com/mtlynch/picoshare/v2/store/test_sqlite"
	"github.com/mtlynch/picoshare/v2/types"
)

func TestGetStorageSpace(t *testing.T) {
	dataStore := test_sqlite.New()

	s := handlers.New(mockAuthenticator{}, dataStore)

	req, err := http.NewRequest("GET", "/api/storage-space", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("GET /api/storage-space returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}


}
