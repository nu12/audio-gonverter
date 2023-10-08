package helper

import (
	"testing"

	"github.com/gorilla/sessions"
	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/database"
	"github.com/nu12/audio-gonverter/internal/user"
)

func TestFlash(t *testing.T) {
	testApp := &config.Config{
		SessionStore: sessions.NewCookieStore([]byte("test")),
		DatabaseRepo: &database.MockDB{Messages: []string{}},
	}

	user := user.New()

	expected := "Test message"

	h := WithConfig(testApp)
	h.AddFlash(user, expected)
	got := h.GetFlash(user)[0]

	if expected != got {
		t.Errorf("Expected %s, got %s", expected, got)
	}
}

func TestSliceToString(t *testing.T) {
	s := []string{"t1", "t2"}
	ps := SliceToString(s)
	if ps != ".t1,.t2" {
		t.Errorf("Expected %s, got %s", ".t1,.t2", ps)
	}
}
