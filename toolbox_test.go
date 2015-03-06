package silkroad

import "testing"

func TestStringInSlice(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if got, want := stringInSlice(slice, "a"), true; got != want {
		t.Errorf("TestStringInSlice got %v, want %v", got, want)
	}

	if got, want := stringInSlice(slice, "d"), false; got != want {
		t.Errorf("TestStringInSlice got %v, want %v", got, want)
	}

}
