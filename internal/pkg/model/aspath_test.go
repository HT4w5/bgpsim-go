package model_test

import (
	"fmt"
	"testing"

	"github.com/HT4w5/bgpsim-go/internal/pkg/model"
)

func TestAsPath_Prepend(t *testing.T) {
	tests := []struct {
		name        string
		initialPath []uint32
		asToAdd     uint32
		expectedErr error
		expectedLen int
		expectedMap map[uint32]int
	}{
		{
			name:        "Add new AS",
			initialPath: []uint32{},
			asToAdd:     1234,
			expectedErr: nil,
			expectedLen: 1,
			expectedMap: map[uint32]int{1234: 0},
		},
		{
			name:        "Add duplicate AS",
			initialPath: []uint32{1234},
			asToAdd:     1234,
			expectedErr: fmt.Errorf("duplicate AS number in path: 1234"),
			expectedLen: 1,
			expectedMap: map[uint32]int{1234: 0},
		},
		{
			name:        "Add multiple AS",
			initialPath: []uint32{2345, 3456},
			asToAdd:     1234,
			expectedErr: nil,
			expectedLen: 3,
			expectedMap: map[uint32]int{1234: 2, 2345: 0, 3456: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := model.NewAsPath()
			for _, as := range tt.initialPath {
				if err := ap.Prepend(as); err != nil {
					t.Fatalf("Unexpected error while setting up test: %v", err)
				}
			}
			err := ap.Prepend(tt.asToAdd)
			if err != tt.expectedErr && (err == nil || tt.expectedErr == nil || err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
			if ap.Len() != tt.expectedLen {
				t.Errorf("Expected path length %d, got %d", tt.expectedLen, ap.Len())
			}
		})
	}
}

func TestAsPath_HasAs(t *testing.T) {
	tests := []struct {
		name        string
		initialPath []uint32
		asToCheck   uint32
		expected    bool
	}{
		{
			name:        "Check AS in empty path",
			initialPath: []uint32{},
			asToCheck:   1234,
			expected:    false,
		},
		{
			name:        "Check existing AS",
			initialPath: []uint32{1234, 2345},
			asToCheck:   1234,
			expected:    true,
		},
		{
			name:        "Check non-existing AS",
			initialPath: []uint32{1234, 2345},
			asToCheck:   3456,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := model.NewAsPath()
			for _, as := range tt.initialPath {
				if err := ap.Prepend(as); err != nil {
					t.Fatalf("Unexpected error while setting up test: %v", err)
				}
			}
			result := ap.HasAs(tt.asToCheck)
			if result != tt.expected {
				t.Errorf("Expected HasAs to return %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAsPath_Len(t *testing.T) {
	tests := []struct {
		name        string
		initialPath []uint32
		expectedLen int
	}{
		{
			name:        "Length of empty path",
			initialPath: []uint32{},
			expectedLen: 0,
		},
		{
			name:        "Length of path with one AS",
			initialPath: []uint32{1234},
			expectedLen: 1,
		},
		{
			name:        "Length of path with multiple AS",
			initialPath: []uint32{1234, 2345, 3456},
			expectedLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := model.NewAsPath()
			for _, as := range tt.initialPath {
				if err := ap.Prepend(as); err != nil {
					t.Fatalf("Unexpected error while setting up test: %v", err)
				}
			}
			result := ap.Len()
			if result != tt.expectedLen {
				t.Errorf("Expected Len to return %d, got %d", tt.expectedLen, result)
			}
		})
	}
}

func TestAsPath_String(t *testing.T) {
	tests := []struct {
		name        string
		initialPath []uint32
		expectedStr string
	}{
		{
			name:        "String representation of empty path",
			initialPath: []uint32{},
			expectedStr: "",
		},
		{
			name:        "String representation of path with one AS",
			initialPath: []uint32{1234},
			expectedStr: "1234",
		},
		{
			name:        "String representation of path with multiple AS",
			initialPath: []uint32{1234, 2345, 3456},
			expectedStr: "3456 2345 1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := model.NewAsPath()
			for _, as := range tt.initialPath {
				if err := ap.Prepend(as); err != nil {
					t.Fatalf("Unexpected error while setting up test: %v", err)
				}
			}
			result := ap.String()
			if result != tt.expectedStr {
				t.Errorf("Expected String to return '%s', got '%s'", tt.expectedStr, result)
			}
		})
	}
}

func equalMaps(a, b map[uint32]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
