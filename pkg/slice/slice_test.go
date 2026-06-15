package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedupNoDuplicates(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := Dedup(input)
	assert.Equal(t, input, got, "slice with no duplicates should be returned unchanged")
}

func TestDedupWithDuplicates(t *testing.T) {
	input := []int{1, 2, 2, 3, 4, 3, 5}
	want := []int{1, 2, 3, 4, 5}
	got := Dedup(input)
	assert.Equal(t, want, got, "duplicates should be removed preserving order of first occurrence")
}

func TestDedupAllDuplicates(t *testing.T) {
	input := []int{7, 7, 7, 7}
	want := []int{7}
	got := Dedup(input)
	assert.Equal(t, want, got, "all duplicate elements should return a single element")
}

func TestDedupEmpty(t *testing.T) {
	var input []int
	got := Dedup(input)
	assert.Empty(t, got, "empty slice should return empty")
	assert.Nil(t, got, "empty slice should return nil")
}

func TestDedupString(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b", "d"}
	want := []string{"a", "b", "c", "d"}
	got := Dedup(input)
	assert.Equal(t, want, got, "string duplicates should be removed preserving order")
}

func TestDedupInt(t *testing.T) {
	input := []int{3, 1, 2, 3, 1, 4}
	want := []int{3, 1, 2, 4}
	got := Dedup(input)
	assert.Equal(t, want, got, "int duplicates should be removed preserving order")
}
