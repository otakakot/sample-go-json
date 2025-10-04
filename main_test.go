package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	goccyjson "github.com/goccy/go-json"
)

func HandleJSONMarshal(w http.ResponseWriter, r *http.Request) {
	var req map[string]any

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)

		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(req)
}

func HandleJSONDecode(w http.ResponseWriter, r *http.Request) {
	var req map[string]any

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(req)
}

func HandleJSONPipe(w http.ResponseWriter, r *http.Request) {
	var req map[string]any

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		_, err := io.Copy(pw, r.Body)
		if err != nil {
			pw.CloseWithError(err)
		}
	}()

	if err := json.NewDecoder(pr).Decode(&req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

func HandleGoccyJSONUnmarshal(w http.ResponseWriter, r *http.Request) {
	var req map[string]any

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err := goccyjson.Unmarshal(body, &req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	goccyjson.NewEncoder(w).Encode(req)
}

func HandleGoccyJSONDecode(w http.ResponseWriter, r *http.Request) {
	var req map[string]any

	if err := goccyjson.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	goccyjson.NewEncoder(w).Encode(req)
}

func HandleGoccyJSONPipe(w http.ResponseWriter, r *http.Request) {
	var req map[string]any

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		_, err := io.Copy(pw, r.Body)
		if err != nil {
			pw.CloseWithError(err)
		}
	}()

	if err := goccyjson.NewDecoder(pr).Decode(&req); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	goccyjson.NewEncoder(w).Encode(req)
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	jsonData := []byte(`{
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "Tokyo",
			"country": "Japan"
		},
		"tags": ["developer", "golang", "backend"]
	}`)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleJSONMarshal(w, req)
	}
}

func BenchmarkJSONDecode(b *testing.B) {
	jsonData := []byte(`{
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "Tokyo",
			"country": "Japan"
		},
		"tags": ["developer", "golang", "backend"]
	}`)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleJSONDecode(w, req)
	}
}

func BenchmarkJSONPipe(b *testing.B) {
	jsonData := []byte(`{
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "Tokyo",
			"country": "Japan"
		},
		"tags": ["developer", "golang", "backend"]
	}`)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleJSONPipe(w, req)
	}
}

func BenchmarkGoccyJSONUnmarshal(b *testing.B) {
	jsonData := []byte(`{
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "Tokyo",
			"country": "Japan"
		},
		"tags": ["developer", "golang", "backend"]
	}`)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleGoccyJSONUnmarshal(w, req)
	}
}

func BenchmarkGoccyJSONDecode(b *testing.B) {
	jsonData := []byte(`{
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "Tokyo",
			"country": "Japan"
		},
		"tags": ["developer", "golang", "backend"]
	}`)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleGoccyJSONDecode(w, req)
	}
}

func BenchmarkGoccyJSONPipe(b *testing.B) {
	jsonData := []byte(`{
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "Tokyo",
			"country": "Japan"
		},
		"tags": ["developer", "golang", "backend"]
	}`)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleGoccyJSONPipe(w, req)
	}
}

func generateLargeJSON() []byte {
	var sb strings.Builder
	sb.WriteString(`{"users":[`)

	for i := range 500 {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`{
			"id": %d,
			"name": "User %d",
			"email": "user%d@example.com",
			"age": %d,
			"address": {
				"street": "%d Main Street",
				"city": "City %d",
				"country": "Country %d",
				"zipcode": "%05d"
			},
			"tags": ["tag1", "tag2", "tag3", "tag4", "tag5"],
			"metadata": {
				"created_at": "2024-01-01T00:00:00Z",
				"updated_at": "2024-12-31T23:59:59Z",
				"status": "active",
				"verified": true
			}
		}`, i, i, i, 20+i%50, i, i, i, i))
	}
	sb.WriteString(`]}`)

	return []byte(sb.String())
}

func BenchmarkJSONUnmarshalLargePayload(b *testing.B) {
	jsonData := generateLargeJSON()
	b.Logf("Large JSON payload size: %d bytes (%.2f KB)", len(jsonData), float64(len(jsonData))/1024)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleJSONMarshal(w, req)
	}
}

func BenchmarkJSONDecodeLargePayload(b *testing.B) {
	jsonData := generateLargeJSON()
	b.Logf("Large JSON payload size: %d bytes (%.2f KB)", len(jsonData), float64(len(jsonData))/1024)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleJSONDecode(w, req)
	}
}

func BenchmarkJSONPipeLargePayload(b *testing.B) {
	jsonData := generateLargeJSON()
	b.Logf("Large JSON payload size: %d bytes (%.2f KB)", len(jsonData), float64(len(jsonData))/1024)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleJSONPipe(w, req)
	}
}

func BenchmarkGoccyJSONUnmarshalLargePayload(b *testing.B) {
	jsonData := generateLargeJSON()
	b.Logf("Large JSON payload size: %d bytes (%.2f KB)", len(jsonData), float64(len(jsonData))/1024)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleGoccyJSONUnmarshal(w, req)
	}
}

func BenchmarkGoccyJSONDecodeLargePayload(b *testing.B) {
	jsonData := generateLargeJSON()
	b.Logf("Large JSON payload size: %d bytes (%.2f KB)", len(jsonData), float64(len(jsonData))/1024)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleGoccyJSONDecode(w, req)
	}
}

func BenchmarkGoccyJSONPipeLargePayload(b *testing.B) {
	jsonData := generateLargeJSON()
	b.Logf("Large JSON payload size: %d bytes (%.2f KB)", len(jsonData), float64(len(jsonData))/1024)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonData))
		w := httptest.NewRecorder()
		HandleGoccyJSONPipe(w, req)
	}
}
