package ionCannonClient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIonCannonClient_CheckStatus(t *testing.T) {
	// Create a mock HTTP server to simulate the Ion Cannon API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status" {
			t.Errorf("Expected request to /status, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"generation": 1, "available": true}`))
	}))
	defer server.Close()

	client := NewIonCannonClient(server.URL)
	status, err := client.CheckStatus()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !status.Available {
		t.Errorf("Expected Ion Cannon to be available, got false")
	}
}

func TestIonCannonClient_FireCommand(t *testing.T) {
	// Create a mock HTTP server to simulate the Ion Cannon API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/fire" {
			t.Errorf("Expected request to /fire, got %s", r.URL.Path)
		}

		var requestBody struct {
			Target  map[string]int `json:"target"`
			Enemies int            `json:"enemies"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to parse request body: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"casualties": 1, "generation": %d}`, requestBody.Enemies)))
	}))
	defer server.Close()

	client := NewIonCannonClient(server.URL)
	casualties, generation, err := client.FireCommand(0, 40, 1)

	// Check the results
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if casualties != 1 {
		t.Errorf("Expected casualties to be 1, got %d", casualties)
	}
	if generation != 1 {
		t.Errorf("Expected generation to be 1, got %d", generation)
	}
}
