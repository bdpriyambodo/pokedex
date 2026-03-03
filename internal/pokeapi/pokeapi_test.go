package pokeapi

import (
	"reflect"
	"testing"
)

// func TestPokeApiRaw(t *testing.T) {
// 	cases := []struct{
// 		name        string
// 		input       []byte
// 		expected    []string
// 		expectError bool
// 	}{
// 		input:=
// 	}
// }

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
		actual, err := PokeApiData(tc.input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("expected %v, got %v", tc.expected, actual)
		}
	}
}
