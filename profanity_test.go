package main

import "testing"

func TestCleanChirp(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "this is a clean string",
			expected: "this is a clean string",
		}, {
			input:    "this is definitely shArbert not a clean string",
			expected: "this is definitely **** not a clean string",
		}, {
			input:    "Fornax me in my cute little kerfuFFle",
			expected: "**** me in my cute little ****",
		}, {
			input:    "We were fornax-ing all sharbert night",
			expected: "We were fornax-ing all **** night",
		}, {
			input:    "Kerfuffle kerfuffle. Sharbert sharbert! Fornax fornax?",
			expected: "**** kerfuffle. **** sharbert! **** fornax?",
		}, {
			input:    "",
			expected: "",
		},
	}

	for _, c := range cases {
		actual := cleanChirp(c.input)

		if actual != c.expected {
			t.Errorf("Fail: did not properly censor the message. \nActual: %s\nExpected: %s", actual, c.expected)
		}
	}
}
