package httprouterserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandlerPositive(t *testing.T) {
	jsonStr := []byte("{\"a\":5, \"b\":8}")

	request, err := http.NewRequest("PUT", "url", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Failed to create a request: %s", err)
		t.FailNow()
	}

	response := httptest.NewRecorder()

	handler := NewRESTHandler(nil)
	handler.CalculateHandler(response, request, nil)

	resp := response.Result()
	var req factRequest
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		t.Errorf("Failed to decode a body: %s", err)
	}
	if (req.A != 120) || (req.B != 40320) {
		t.Errorf("a must be 120, b must be 40320 but a: %d, b: %d", req.A, req.B)
	}

	jsonStr = []byte("{\"a\":10, \"b\":12}")

	request, err = http.NewRequest("PUT", "url", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Failed to create a request: %s", err)
	}

	response = httptest.NewRecorder()
	handler.CalculateHandler(response, request, nil)

	resp = response.Result()
	req = factRequest{0, 0}
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		t.Errorf("Failed to decode a body: %s", err)
	}
	if (req.A != 3628800) || (req.B != 479001600) {
		t.Errorf("a must be 3628800, b must be 479001600 but a: %d, b: %d", req.A, req.B)
	}

	// test cached values
	request, err = http.NewRequest("PUT", "url", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Failed to create a request: %s", err)
	}
	response = httptest.NewRecorder()
	handler.CalculateHandler(response, request, nil)

	resp = response.Result()
	req = factRequest{0, 0}
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		t.Errorf("Failed to decode a body: %s", err)
	}
	if (req.A != 3628800) || (req.B != 479001600) {
		t.Errorf("a must be 3628800, b must be 479001600 but a: %d, b: %d", req.A, req.B)
	}
}

func TestCalculateHandlerWrongInput(t *testing.T) {
	jsonStr := []byte("{\"a\":0, \"b\":8}")

	request, err := http.NewRequest("PUT", "url", bytes.NewBuffer(jsonStr))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Errorf("Failed to create a request: %s", err)
	}

	response := httptest.NewRecorder()

	handler := NewRESTHandler(nil)
	handler.CalculateHandler(response, request, nil)

	resp := response.Result()
	if resp.StatusCode != 400 {
		t.Error("Status code must be 400")
	}

	var req errorResponse
	err = json.NewDecoder(resp.Body).Decode(&req)
	if err != nil {
		t.Errorf("Failed to parse json: %s", err)
	}
	if req.Error != "Incorrect input" {
		t.Errorf("Expected error string \"Incorrect input\" but returned: %s", req.Error)
	}
}
