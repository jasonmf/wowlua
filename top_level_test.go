package wowlua

import (
	"testing"
)

func TestMultipleTopLevelValues(t *testing.T) {
	tab, err := ParseLua(sample_data)
	if err != nil {
		t.Errorf("Unexpected error parsing data: %q", err)
	}
	t.Logf("got keys: %v", tab.Keys())
	expected_length := 2
	length := tab.Len()
	if length != expected_length {
		t.Errorf("Expected table to have %v entries, got %v", expected_length, length)
	}
}
