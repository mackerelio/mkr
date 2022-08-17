package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

func init() {
	// Shorten interval in tests
	minInterval = 1 * time.Second
}

func TestRequestWithRetry_Success(t *testing.T) {
	var counter int
	var err error
	f0 := func() error {
		counter++
		return nil
	}

	counter = 0
	err = requestWithRetry(f0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if counter != 1 {
		t.Errorf("function should be called only once, but called %d times", counter)
	}

	counter = 0
	err = requestWithRetry(f0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if counter != 1 {
		t.Errorf("function should be called only once, but called %d times", counter)
	}

	counter = 0
	err = requestWithRetry(f0, 5)
	if err != nil {
		t.Fatal(err)
	}
	if counter != 1 {
		t.Errorf("function should be called only once, but called %d times", counter)
	}
}

func TestRequestWithRetry_Giveup(t *testing.T) {
	var counter int
	var err error
	f0 := func() error {
		counter++
		return fmt.Errorf("ohno")
	}

	counter = 0
	err = requestWithRetry(f0, 1)
	if err == nil {
		t.Error("error should occur")
	}
	if counter != 2 { // 1 + 1 retry
		t.Errorf("function should be called 2 times, but called %d times", counter)
	}

	counter = 0
	err = requestWithRetry(f0, 0)
	if err == nil {
		t.Error("error should occur")
	}
	if counter != 1 { // 1 + 0 retry
		t.Errorf("function should be called only once, but called %d times", counter)
	}

	counter = 0
	err = requestWithRetry(f0, 5)
	if err == nil {
		t.Error("error should occur")
	}
	if counter != 6 { // 1 + 5 retry
		t.Errorf("function should be called 6 times, but called %d times", counter)
	}
}

func TestRequestWithRetry_Recovery(t *testing.T) {
	var counter int
	var err error
	f0 := func() error {
		counter++
		if counter < 3 {
			return fmt.Errorf("Not yet")
		}
		return nil
	}

	counter = 0
	err = requestWithRetry(f0, 1)
	if err == nil {
		t.Error("error should occur")
	}
	if counter != 2 { // 1 + 1 retry
		t.Errorf("function should be called 2 times, but called %d times", counter)
	}

	counter = 0
	err = requestWithRetry(f0, 5)
	if err != nil {
		t.Error("error should occur")
	}
	if counter != 3 { // Success on 3rd try
		t.Errorf("function should be called 3 times, but called %d times", counter)
	}
}

func TestRequestWithRetry_Status(t *testing.T) {
	var counter int
	var err error
	var status int
	f0 := func() error {
		counter++
		return &mackerel.APIError{StatusCode: status, Message: "ohno"}
	}

	counter = 0
	status = 500
	err = requestWithRetry(f0, 1)
	if err == nil {
		t.Error("error should occur")
	}
	if counter != 2 { // 1 + 1 retry
		t.Errorf("function should be called 2 times, but called %d times", counter)
	}

	counter = 0
	status = 403
	err = requestWithRetry(f0, 1)
	if err == nil {
		t.Error("error should occur")
	}
	if counter != 1 { // 1 + 0 retry
		t.Errorf("function should be called only once, but called %d times", counter)
	}
}
