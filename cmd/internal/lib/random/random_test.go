package random

import "testing"

func TestRandomStringLength(t *testing.T) {
	var tests = []struct {
		name  string
		input int
		want  int
	}{
		{"Length is 0", 0, 4},
		{"Length is 1", 1, 4},
		{"Length is 2", 2, 4},
		{"Length is 3", 3, 4},
		{"Length is 10", 10, 10},
		{"Length is 20", 20, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := len(RandomString(tt.input))
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestRandomStringRandomness(t *testing.T) {
	res1 := RandomString(5)
	res2 := RandomString(5)
	if res1 == res2 {
		t.Errorf("Result incorrect. Got %s for res1 and %s for res2", res1, res2)
	}
}