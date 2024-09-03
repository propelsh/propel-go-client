package propel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetUserByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/users/12345" {
            t.Fatalf("expected to request '/users/12345', got: %s", r.URL.Path)
        }

        // Simulate a successful response
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(User{
            ID:    "12345",
            Email: "john.doe@example.com",
			CreatedAt: time.Now(),
        })
    }))
    defer server.Close()

    client := NewClient(server.URL, "test-api-key")

    // Call the method
    user, err := client.GetUserByID("12345")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Validate the response
    expected := &User{
        ID:    "12345",
        Email: "john.doe@example.com",
		CreatedAt: user.CreatedAt,
    }
    if !reflect.DeepEqual(user, expected) {
        t.Errorf("expected %+v, got %+v", expected, user)
    }
}
