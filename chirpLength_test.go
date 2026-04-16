package main

import "testing"

func TestChirpLength(t *testing.T) {
	cases := []struct {
		input        string
		expectedBool bool
		expectedMSG  string
	}{
		{
			input:        "This is a short chirp",
			expectedBool: true,
			expectedMSG:  "This is a short chirp",
		},
		{
			input:        "This is a ssssssssssssssssssssssssssuuuuuuuuuuppppppppppppppeeeeeeeeeeerrrrrrrrrrr lllllllllllooooooooooooooonnnnnnnnnnngggggggggggg cccccccccccccccccchhhhhhhhhhhhhhhhhhhhiiiiiiiiiiiiiiiiiiiiiiirrrrrrrrrrrrrrrrrrrrppppppppppppppppppppppp",
			expectedBool: false,
			expectedMSG:  "Chirp is too long",
		},
	}

	for _, c := range cases {
		actualBool, actualMSG := validateChirp(c.input)

		if actualBool != c.expectedBool || actualMSG != c.expectedMSG {
			t.Errorf("Fail: Length not checked properly. \nActual: %t,%s\nExpected: %t, %s", actualBool, actualMSG, c.expectedBool, c.expectedMSG)
		}
	}
}
