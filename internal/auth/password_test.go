package auth

import "testing"

func TestPasswordHash(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "password",
			expected: true,
		}, {
			input:    "sup3rSe8crtP4ssw0rd",
			expected: true,
		}, {
			input:    "special)(^(*&$^&%#))characters",
			expected: true,
		}, {
			input:    "l33tsp35k",
			expected: true,
		},
	}

	for _, c := range cases {
		hash, err1 := HashPassword(c.input)

		actual, err2 := CheckPasswordHash(c.input, hash)

		if actual != c.expected || err1 != nil || err2 != nil {
			t.Errorf("Fail: hash incorrect.\nActual: %t\nExpected: %t\nErr1: %s	Err2: %s", actual, c.expected, err1, err2)
		}
	}
}
