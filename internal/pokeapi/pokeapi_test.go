package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPokeApiRaw(t *testing.T) {
	cases := []struct {
		name           string
		offset         int
		responseBody   string
		responseStatus int
		expectError    bool
	}{
		{
			name:           "success",
			offset:         0,
			responseBody:   `{"results":[{"name":"area1"}]}`,
			responseStatus: 200,
			expectError:    false,
		},
		{
			name:           "server error",
			offset:         0,
			responseBody:   `internal error`,
			responseStatus: 500,
			expectError:    false, // depends if you handle status code
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			// Create fake test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.responseStatus)
				w.Write([]byte(tc.responseBody))
			}))
			defer server.Close()

			client := server.Client()

			body, err := PokeApiRaw(client, server.URL)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(body) != tc.responseBody {
				t.Errorf("expected %s, got %s", tc.responseBody, string(body))
			}
		})
	}
}

func TestPokeApiData(t *testing.T) {
	cases := []struct {
		name        string
		input       []byte
		expected    []string
		expectError bool
	}{
		{
			name: "valid multiple areas",
			input: []byte(`{
				"results": [
					{"name": "canalave-city-area"},
					{"name": "eterna-city-area"},
					{"name": "pastoria-city-area"}
				]
			}`),
			expected: []string{
				"canalave-city-area",
				"eterna-city-area",
				"pastoria-city-area",
			},
			expectError: false,
		},
		{
			name: "empty results",
			input: []byte(`{
				"results": []
			}`),
			expected:    []string{},
			expectError: false,
		},
	}

	for _, tc := range cases {
		actual, err := PokeApiLocationArea(tc.input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("expected %v, got %v", tc.expected, actual)
		}
	}
}
