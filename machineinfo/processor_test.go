package machineinfo

import "testing"

func TestMem(t *testing.T) {}
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"", ""}
	}
	for _, c := range cases {
		got := Reverse(c.in)
		if got != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
