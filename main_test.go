package main

import (
	"testing"
)

func TestDedupe(t *testing.T) {
	left := make(chan uint, 10)
	right := make(chan uint, 10)
	result := make(chan uint)

	for _, v := range []uint{1, 2, 5, 6} {
		left <- v
	}
	close(left)

	for _, v := range []uint{2, 4, 5, 7} {
		right <- v
	}
	close(right)

	go Dedupe(result, left, right)
	expected := []uint{1, 2, 4, 5, 6, 7}
	for _, v := range expected {
		next := <-result

		if v != next {
			t.Errorf("Expected %v, got %v", v, next)
		}
	}
	if _, ok := <-result; ok {
		t.Error("Expected result channel to be closed")
	}
}

func TestCapAt(t *testing.T) {
	input := make(chan uint, 6)
	result := make(chan uint)

	for _, v := range []uint{2, 4, 6, 8, 10, 12} {
		input <- v
	}

	go CapAt(10, result, input)
	expected := []uint{2, 4, 6, 8}

	for _, v := range expected {
		next := <-result

		if v != next {
			t.Errorf("Expected %v, got %v", v, next)
		}
	}
	if _, ok := <-result; ok {
		t.Error("Expected result channel to be closed")
	}
}
