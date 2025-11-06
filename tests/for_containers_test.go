package tests

import (
	"net/http"
	"testing"
	"time"
)

func TestPagesAvailability(t *testing.T) {
	baseURL := "http://localhost:7540"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	t.Run("MainPage", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/")
		if err != nil {
			t.Fatalf("GET / failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET / returned status %d, expected 200", resp.StatusCode)
		} else {
			t.Logf("✓ GET / - Status: 200 OK")
		}
	})

	t.Run("LoginPage", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/login.html")
		if err != nil {
			t.Fatalf("GET /login.html failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /login.html returned status %d, expected 200", resp.StatusCode)
		} else {
			t.Logf("✓ GET /login.html - Status: 200 OK")
		}
	})
}