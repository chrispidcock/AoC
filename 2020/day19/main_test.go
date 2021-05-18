package main

import (
	"testing"
)

func TestListRecursive(t *testing.T) {
	i := 0
	matches := [][]string{{"a", "b", "c"}, {"d", "e"}, {"f"}}
	want := []string{"adf", "aef", "bdf", "bef", "cdf", "cef"}
	results := ListRecursive(i, matches)
	for i := range want {
		if want[i] != results[i] {
			t.Fatal(results, want, results[i], want[i])
		}
	}
	if len(results) != len(want) {
		t.Fatal("Lengths do not match: ", "(want)", len(want), "!=", len(results), "(results)")
	}
}

func TestGenMatches(t *testing.T) {
	i := 0
	matches := [][]string{{"a", "b", "c"}, {"d", "e"}, {"f"}}
	want := []string{"adf", "aef", "bdf", "bef", "cdf", "cef"}
	results := ListRecursive(i, matches)
	for i := range want {
		if want[i] != results[i] {
			t.Fatal(want, results, want[i], results[i])
		}
	}
	if len(results) != len(want) {
		t.Fatal("Lengths do not match: ", "(want)", len(want), "!=", len(results), "(results)")
	}
}
