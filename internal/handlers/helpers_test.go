package handlers

import (
	"testing"
)

func TestFormatJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		hasError bool
	}{
		{
			name:  "simple map",
			input: map[string]interface{}{"key": "value"},
			expected: `{
  "key": "value"
}`,
			hasError: false,
		},
		{
			name:  "nested map",
			input: map[string]interface{}{"data": map[string]interface{}{"id": "123", "name": "test"}},
			expected: `{
  "data": {
    "id": "123",
    "name": "test"
  }
}`,
			hasError: false,
		},
		{
			name:  "slice",
			input: []string{"item1", "item2"},
			expected: `[
  "item1",
  "item2"
]`,
			hasError: false,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "null",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatJSONResponse(tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestGetStringArg(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		key      string
		expected string
	}{
		{
			name:     "valid string",
			args:     map[string]interface{}{"key": "value"},
			key:      "key",
			expected: "value",
		},
		{
			name:     "missing key",
			args:     map[string]interface{}{"other": "value"},
			key:      "key",
			expected: "",
		},
		{
			name:     "empty string",
			args:     map[string]interface{}{"key": ""},
			key:      "key",
			expected: "",
		},
		{
			name:     "wrong type",
			args:     map[string]interface{}{"key": 123},
			key:      "key",
			expected: "",
		},
		{
			name:     "nil value",
			args:     map[string]interface{}{"key": nil},
			key:      "key",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStringArg(tt.args, tt.key)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetIntArg(t *testing.T) {
	tests := []struct {
		name         string
		args         map[string]interface{}
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid int",
			args:         map[string]interface{}{"key": 42.0},
			key:          "key",
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "missing key",
			args:         map[string]interface{}{"other": 42.0},
			key:          "key",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "wrong type",
			args:         map[string]interface{}{"key": "not a number"},
			key:          "key",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "zero value",
			args:         map[string]interface{}{"key": 0.0},
			key:          "key",
			defaultValue: 10,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetIntArg(tt.args, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetStringArrayArg(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		key      string
		expected []string
	}{
		{
			name:     "valid string array",
			args:     map[string]interface{}{"key": []interface{}{"item1", "item2"}},
			key:      "key",
			expected: []string{"item1", "item2"},
		},
		{
			name:     "empty array",
			args:     map[string]interface{}{"key": []interface{}{}},
			key:      "key",
			expected: []string{},
		},
		{
			name:     "missing key",
			args:     map[string]interface{}{"other": []interface{}{"item1"}},
			key:      "key",
			expected: []string{},
		},
		{
			name:     "wrong type",
			args:     map[string]interface{}{"key": "not an array"},
			key:      "key",
			expected: []string{},
		},
		{
			name:     "mixed types in array",
			args:     map[string]interface{}{"key": []interface{}{"item1", 123, "item2"}},
			key:      "key",
			expected: []string{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStringArrayArg(tt.args, tt.key)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("Expected %q at index %d, got %q", tt.expected[i], i, v)
				}
			}
		})
	}
}

func TestCreateSimpleResponse(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		message  string
		expected map[string]interface{}
	}{
		{
			name:    "with message",
			data:    []string{"item1", "item2"},
			message: "Success",
			expected: map[string]interface{}{
				"data":    []string{"item1", "item2"},
				"count":   2,
				"message": "Success",
			},
		},
		{
			name:    "without message",
			data:    []string{"item1"},
			message: "",
			expected: map[string]interface{}{
				"data":  []string{"item1"},
				"count": 1,
			},
		},
		{
			name:    "non-slice data",
			data:    map[string]string{"key": "value"},
			message: "Test",
			expected: map[string]interface{}{
				"data":    map[string]string{"key": "value"},
				"message": "Test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateSimpleResponse(tt.data, tt.message)

			// Check data (handle slices specially)
			if expectedSlice, ok := tt.expected["data"].([]string); ok {
				if resultSlice, ok := result["data"].([]string); ok {
					if len(expectedSlice) != len(resultSlice) {
						t.Errorf("Expected slice length %d, got %d", len(expectedSlice), len(resultSlice))
					} else {
						for i, v := range expectedSlice {
							if v != resultSlice[i] {
								t.Errorf("Expected slice[%d] %q, got %q", i, v, resultSlice[i])
							}
						}
					}
				} else {
					t.Errorf("Expected slice data, got %T", result["data"])
				}
			} else if expectedMap, ok := tt.expected["data"].(map[string]string); ok {
				if resultMap, ok := result["data"].(map[string]string); ok {
					if len(expectedMap) != len(resultMap) {
						t.Errorf("Expected map length %d, got %d", len(expectedMap), len(resultMap))
					} else {
						for k, v := range expectedMap {
							if resultMap[k] != v {
								t.Errorf("Expected map[%q] %q, got %q", k, v, resultMap[k])
							}
						}
					}
				} else {
					t.Errorf("Expected map data, got %T", result["data"])
				}
			} else {
				if result["data"] != tt.expected["data"] {
					t.Errorf("Expected data %v, got %v", tt.expected["data"], result["data"])
				}
			}

			// Check count (if expected)
			if expectedCount, ok := tt.expected["count"]; ok {
				if result["count"] != expectedCount {
					t.Errorf("Expected count %v, got %v", expectedCount, result["count"])
				}
			}

			// Check message (if expected)
			if expectedMessage, ok := tt.expected["message"]; ok {
				if result["message"] != expectedMessage {
					t.Errorf("Expected message %v, got %v", expectedMessage, result["message"])
				}
			}
		})
	}
}
