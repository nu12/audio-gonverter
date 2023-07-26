package main

import (
	"net/http"
	"os"
	"testing"
	"time"
)

// Test loadEnv
func TestLoadEnv(t *testing.T) {
	testConfig := &Config{
		Env: map[string]string{},
	}

	if err := testConfig.loadEnv([]string{"EMPTY"}); err == nil {
		t.Errorf("Empty env variable should return an error, but error is nil")
	}

	env := "EXISTS"
	val := "yes"
	os.Setenv(env, val)
	if err := testConfig.loadEnv([]string{env}); err != nil {
		t.Errorf("Existing env variable should not return an error, but got %s", err)
	}

	if testConfig.Env[env] != val {
		t.Errorf("Config should contain %s with value %s", env, val)
	}
}

// Test startWeb
type TestServer struct {
	Addr    string
	Handler http.Handler
}

func (*TestServer) ListenAndServe() error {
	return nil
}

func TestStartWeb(t *testing.T) {
	c := make(chan error)
	app := Config{}

	t.Run("Start Web Service", func(t *testing.T) {
		testServer := &TestServer{}
		go app.startWeb(c, testServer)

		select {
		case err := <-c:
			t.Errorf("Unexpected error: %s", err)
		case <-time.After(1 * time.Second):
			// No error occurred, the test passes
		}
	})
}
