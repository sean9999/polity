package subject

import "testing"

func TestSubject_String(t *testing.T) {
	cases := []struct {
		name string
		in   Subject
		want string
	}{
		{"empty", Subject(""), ""},
		{"simple", Subject("hello"), "hello"},
		{"unicode", Subject("héllo 🚀"), "héllo 🚀"},
	}
	for _, tc := range cases {
		// capture range var
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.String()
			if got != tc.want {
				t.Fatalf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Subject
	}{
		{"empty", "", Subject("")},
		{"simple", "hello", Subject("hello")},
		{"unicode", "héllo 🚀", Subject("héllo 🚀")},
	}
	for _, tc := range cases {
		// capture range var
		t.Run(tc.name, func(t *testing.T) {
			got := From(tc.in)
			if got != tc.want {
				t.Fatalf("From(%q) = %v, want %v", tc.in, got, tc.want)
			}
			// also verify round-trip via String()
			if got.String() != tc.in {
				t.Fatalf("From(%q).String() = %q, want %q", tc.in, got.String(), tc.in)
			}
		})
	}
}
