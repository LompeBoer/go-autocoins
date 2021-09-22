package autocoins

import "testing"

func TestContainsString(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e", "g"}
	if !ContainsString(a, "a") {
		t.Errorf("search invalid: should find a")
	}
	if !ContainsString(a, "c") {
		t.Errorf("search invalid: should find c")
	}
	if !ContainsString(a, "g") {
		t.Errorf("search invalid: should find g")
	}
	if ContainsString(a, "f") {
		t.Errorf("search invalid: should not find f")
	}
	if ContainsString(a, "i") {
		t.Errorf("search invalid: should not find i")
	}
}

func TestContainsStringSorted(t *testing.T) {
	a := []string{"a", "b", "c", "d", "e", "g"}
	if !ContainsStringSorted(a, "a") {
		t.Errorf("search invalid: should find a")
	}
	if !ContainsString(a, "c") {
		t.Errorf("search invalid: should find c")
	}
	if !ContainsString(a, "g") {
		t.Errorf("search invalid: should find g")
	}
	if ContainsString(a, "f") {
		t.Errorf("search invalid: should not find f")
	}
	if ContainsString(a, "i") {
		t.Errorf("search invalid: should not find i")
	}
}
