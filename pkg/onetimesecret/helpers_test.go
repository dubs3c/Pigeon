package onetimesecret

import (
	"testing"
	"time"
)

func TestIsSecretValid(t *testing.T) {
	currentTime := time.Now()
	secret := Secret{
		token:    "123",
		message:  "string",
		password: "string",
		expire:   "2019-10-31 10:05:00",
		maxviews: 1,
		views:    0,
	}

	valid := IsSecretValid(currentTime, &secret)

	if valid != false {
		t.Errorf("valid = %t; want false", valid)
	}

	loc, _ := time.LoadLocation("Europe/Berlin")
	currentTime, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-10-30 10:05:00", loc)
	if err != nil {
		t.Errorf("Error parsing time in test")
	}

	valid = IsSecretValid(currentTime, &secret)

	if valid != true {
		t.Errorf("valid = %t; want true", valid)
	}

}
