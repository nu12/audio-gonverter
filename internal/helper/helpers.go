package helper

import (
	"strings"

	"github.com/nu12/audio-gonverter/internal/config"
	"github.com/nu12/audio-gonverter/internal/user"
)

type Helper struct {
	Config *config.Config
}

func WithConfig(c *config.Config) *Helper {
	return &Helper{
		Config: c,
	}
}

func (h *Helper) SaveUser(user *user.User) error {
	return h.Config.DatabaseRepo.Save(user)
}

func (h *Helper) LoadUser(id string) (*user.User, error) {
	return h.Config.DatabaseRepo.Load(id)
}

func (h *Helper) AddFlash(u *user.User, msg string) {
	u.AddMessage(msg)
	if err := h.SaveUser(u); err != nil {
		h.Config.Log.Warning("Error saving user with flash message: " + err.Error())
	}
}

func (h *Helper) GetFlash(u *user.User) []string {
	if len(u.Messages) == 0 {
		return []string{"Welcome to audio-gonverter!"}
	}
	messages := u.GetMessages()
	if err := h.SaveUser(u); err != nil {
		h.Config.Log.Warning("Error saving user without flash messages: " + err.Error())
	}
	return messages
}

func SliceToString(s []string) string {
	ps := []string{}
	for _, f := range s {
		ps = append(ps, "."+f)
	}
	return strings.Join(ps, ",")
}
