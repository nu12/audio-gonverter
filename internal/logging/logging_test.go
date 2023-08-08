package logging

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	println("TEST", "Message")
	expected := createTimeStamp() + " TEST Message\n"
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}
}
func TestDebug(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	os.Setenv("VERBOSITY", "3")
	l := NewLogger()
	l.Debug("Debug message")
	expected := ""
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	os.Setenv("VERBOSITY", "4")
	l = NewLogger()
	l.Debug("message")
	expected = createTimeStamp() + " DEBUG message\n"
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	defer os.Unsetenv("VERBOSITY")
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	os.Setenv("VERBOSITY", "2")
	l := NewLogger()
	l.Info("Debug message")
	expected := ""
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	os.Setenv("VERBOSITY", "3")
	l = NewLogger()
	l.Info("message")
	expected = createTimeStamp() + " INFO message\n"
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	defer os.Unsetenv("VERBOSITY")
}

func TestWarning(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	os.Setenv("VERBOSITY", "1")
	l := NewLogger()
	l.Warning("Debug message")
	expected := ""
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	os.Setenv("VERBOSITY", "2")
	l = NewLogger()
	l.Warning("message")
	expected = createTimeStamp() + " WARNING message\n"
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	defer os.Unsetenv("VERBOSITY")
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	os.Setenv("VERBOSITY", "0")
	l := NewLogger()
	l.Error(errors.New("Debug message"))
	expected := ""
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	os.Setenv("VERBOSITY", "1")
	l = NewLogger()
	l.Error(errors.New("message"))
	expected = createTimeStamp() + " ERROR message\n"
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

	defer os.Unsetenv("VERBOSITY")
}

func TestFatal(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Fatalln did not panic")
		}
	}()

	l := NewLogger()
	l.Fatal(errors.New("message"))
	expected := createTimeStamp() + " FATAL message\n"
	if got := buf.String(); got != expected {
		t.Errorf("Message didn't match. Got: %q. Expected:  %q", got, expected)
	}

}

func createTimeStamp() string {
	now := time.Now()
	timestamp := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	return timestamp
}
