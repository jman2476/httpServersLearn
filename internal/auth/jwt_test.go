package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	dur, _ := time.ParseDuration("3s")
	cases := struct {
		inputID         uuid.UUID
		inputSecret     string
		inputExpiration time.Duration
		expected        string
	}{
		inputID:         uuid.New(),
		inputSecret:     "super secret",
		inputExpiration: dur,
		expected:        "string",
	}

	token, err := MakeJWT(cases.inputID, cases.inputSecret, cases.inputExpiration)
	t.Log(token)
	if err != nil {
		t.Errorf("Fail: error %s making JWT", err)
	}
	if fmt.Sprintf("%T", token) != cases.expected {
		t.Errorf("Fail: token is not of type string, it's %T", token)
	}
}

func TestValidateJWT(t *testing.T) {
	dur5min, _ := time.ParseDuration("5m")
	dur3s, _ := time.ParseDuration("3s")
	dur10h, _ := time.ParseDuration("10h")

	cases := []struct {
		inputID     uuid.UUID
		inputSecret string
		inputExp    time.Duration
		wait        bool
		expected    uuid.UUID
		expectedErr error
	}{
		{
			inputID:     uuid.New(),
			inputSecret: "super duper secret",
			inputExp:    dur5min,
			wait:        false,
			expected:    uuid.UUID{},
			expectedErr: nil,
		}, {
			inputID:     uuid.New(),
			inputSecret: "happy buttons is my name",
			inputExp:    dur5min,
			wait:        true,
			expected:    uuid.UUID{},
			expectedErr: nil,
		}, {
			inputID:     uuid.New(),
			inputSecret: "Why did I eat ice cream on my skateboard?",
			inputExp:    dur3s,
			wait:        true,
			expected:    uuid.UUID{},
			expectedErr: nil,
		}, {
			inputID:     uuid.New(),
			inputSecret: "cheese is melting down my back",
			inputExp:    dur3s,
			wait:        false,
			expected:    uuid.UUID{},
			expectedErr: nil,
		}, {
			inputID:     uuid.New(),
			inputSecret: "donde esta la bibliotecha?",
			inputExp:    dur10h,
			wait:        false,
			expected:    uuid.UUID{},
			expectedErr: nil,
		},
	}

	for _, c := range cases {
		c.expected = c.inputID
		token, err := MakeJWT(c.inputID, c.inputSecret, c.inputExp)
		if err != nil {
			t.Errorf("Duration: %v Wait: %t", c.inputExp, c.wait)
			t.Errorf("Fail: issue making JWT %s", err)
		}

		if c.wait {
			time.Sleep(dur3s)
		}

		id, err := ValidateJWT(token, c.inputSecret)
		if err != nil {
			if c.wait {
				t.Logf("Token expired as expected %v", id)
				continue
			} else {
				t.Errorf("Duration: %v Wait: %t", c.inputExp, c.wait)
				t.Errorf("Fail: token not validated %s", err)
			}
		}

		log.Printf("Type of actual ID: %T", id)
		if id != c.expected {
			t.Errorf("Duration: %v Wait: %t", c.inputExp, c.wait)
			t.Errorf("Fail: id %s does not match expected %s; err %s", id.String(), c.expected.String(), err)
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	id := uuid.New()
	secret := "Super testing secret"
	duration, _ := time.ParseDuration("10m")
	tokenStr, err := MakeJWT(id, secret, duration)
	if err != nil {
		t.Errorf("Fail: Error making token %s", err)
	}

	var headerGood = http.Header{}
	var headerEmpty = http.Header{}
	var headerMalformed = http.Header{}

	val := fmt.Sprintf("Bearer %s", tokenStr)
	headerGood.Add("Authorization", val)
	headerMalformed.Add("Authorization", tokenStr)

	cases := []struct {
		input       http.Header
		expectedStr string
		expectedErr error
	}{
		{
			input:       headerGood,
			expectedStr: tokenStr,
			expectedErr: nil,
		}, {
			input:       headerEmpty,
			expectedStr: "",
			expectedErr: missingBearerTokenErr,
		}, {
			input:       headerMalformed,
			expectedStr: "",
			expectedErr: malformedHeaderErr,
		},
	}

	for _, c := range cases {
		actualStr, actualErr := GetBearerToken(c.input)
		if actualStr != c.expectedStr || !errors.Is(actualErr, c.expectedErr) {
			t.Errorf("Fail: returned token and error mismatch:\nActual: %s | %s \nExpected: %s | %s\n %t %t", actualStr, actualErr, c.expectedStr, c.expectedErr, actualStr != c.expectedStr, actualErr != c.expectedErr)
		}
	}
}
