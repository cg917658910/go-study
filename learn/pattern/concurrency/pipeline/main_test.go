package main

import "testing"

func TestPipeline(t *testing.T) {
	c := generator(2, 3, 4)
	out := square(c)

	expected := []int{4, 9, 16}
	i := 0
	for result := range out {
		if result != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], result)
		}
		i++
	}
}
