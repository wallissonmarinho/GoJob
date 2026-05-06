package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	adapthttp "github.com/wallissonmarinho/GoJob/internal/adapters/http"
)

func TestSyncClient_TriggerSync_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		assert.Equal(t, "POST", r.Method)
		assert.NotEmpty(t, r.Header.Get("Authorization"))
		assert.NotEmpty(t, r.Header.Get("X-Admin-Key"))

		w.WriteHeader(http.StatusAccepted) // 202
		w.Write([]byte(`{"accepted":true}`))
	}))
	defer server.Close()

	client := adapthttp.NewSyncClient(10*time.Second, false)
	statusCode, err := client.TriggerSync(context.Background(), server.URL, "test-key")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, statusCode)
}

func TestSyncClient_TriggerSync_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // 500
		w.Write([]byte(`{"error":"server error"}`))
	}))
	defer server.Close()

	client := adapthttp.NewSyncClient(10*time.Second, false)
	statusCode, err := client.TriggerSync(context.Background(), server.URL, "test-key")

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestSyncClient_TriggerSync_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow server
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with very short timeout
	client := adapthttp.NewSyncClient(100*time.Millisecond, false)
	_, err := client.TriggerSync(context.Background(), server.URL, "test-key")

	assert.Error(t, err)
}

func TestSyncClient_TriggerSync_WithVerbose(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"accepted":true}`))
	}))
	defer server.Close()

	client := adapthttp.NewSyncClient(10*time.Second, true) // verbose enabled
	statusCode, err := client.TriggerSync(context.Background(), server.URL, "test-key")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, statusCode)
}
