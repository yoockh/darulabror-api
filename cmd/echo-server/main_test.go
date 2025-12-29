package main

import (
	"reflect"
	"testing"
)

func TestParseCORSOrigins(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single origin",
			input:    "https://www.darulabror.com",
			expected: []string{"https://www.darulabror.com"},
		},
		{
			name:     "multiple origins",
			input:    "https://www.darulabror.com,https://admin.darulabror.com",
			expected: []string{"https://www.darulabror.com", "https://admin.darulabror.com"},
		},
		{
			name:     "origins with spaces",
			input:    "https://www.darulabror.com, https://admin.darulabror.com , https://api.darulabror.com",
			expected: []string{"https://www.darulabror.com", "https://admin.darulabror.com", "https://api.darulabror.com"},
		},
		{
			name:     "origins with empty entries",
			input:    "https://www.darulabror.com,,https://admin.darulabror.com,",
			expected: []string{"https://www.darulabror.com", "https://admin.darulabror.com"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only commas and spaces",
			input:    " , , , ",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCORSOrigins(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseCORSOrigins(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAppendUniqueOrigins(t *testing.T) {
	tests := []struct {
		name       string
		existing   []string
		newOrigins []string
		expected   []string
	}{
		{
			name:       "no duplicates",
			existing:   []string{"https://www.darulabror.com"},
			newOrigins: []string{"http://localhost:3000"},
			expected:   []string{"https://www.darulabror.com", "http://localhost:3000"},
		},
		{
			name:       "with duplicates",
			existing:   []string{"https://www.darulabror.com", "http://localhost:3000"},
			newOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			expected:   []string{"https://www.darulabror.com", "http://localhost:3000", "http://127.0.0.1:3000"},
		},
		{
			name:       "all duplicates",
			existing:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			newOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
			expected:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		},
		{
			name:       "empty existing",
			existing:   []string{},
			newOrigins: []string{"http://localhost:3000"},
			expected:   []string{"http://localhost:3000"},
		},
		{
			name:       "empty new origins",
			existing:   []string{"https://www.darulabror.com"},
			newOrigins: []string{},
			expected:   []string{"https://www.darulabror.com"},
		},
		{
			name:       "multiple new origins with one duplicate",
			existing:   []string{"https://www.darulabror.com", "https://admin.darulabror.com"},
			newOrigins: []string{"https://admin.darulabror.com", "http://localhost:3000", "http://127.0.0.1:3000"},
			expected:   []string{"https://www.darulabror.com", "https://admin.darulabror.com", "http://localhost:3000", "http://127.0.0.1:3000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := appendUniqueOrigins(tt.existing, tt.newOrigins)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("appendUniqueOrigins(%v, %v) = %v, want %v", tt.existing, tt.newOrigins, result, tt.expected)
			}
		})
	}
}
