package helper

import (
	"fmt"

	"github.com/nu12/audio-gonverter/internal/user"
)

func (h *Helper) StartWorker(c chan<- error) {
	h.Log.Info("Starting Worker service")

	for {
		msg, err := h.Config.QueueRepo.Pull()
		if err != nil {
			c <- err
		}
		decoded, err := h.Config.QueueRepo.Decode(msg)
		if err != nil {
			h.Log.Warning("Cannot decode the message: " + msg)
			continue
		}
		user, err := h.LoadUser(decoded.UserUUID)
		if err != nil {
			h.Log.Warning("Cannot retrieve user: " + decoded.UserUUID)
			continue
		}
		if err := h.Convert(user, decoded.Format, decoded.Kbps); err != nil {
			h.Log.Warning("Error converting file")
		}
		user.IsConverting = false
		if err := h.SaveUser(user); err != nil {
			h.Log.Warning("Error saving user")
			continue
		}
	}
}

func (h *Helper) Convert(user *user.User, format, kpbs string) error {
	for _, file := range user.Files {

		err := h.Config.ConvertionToolRepo.Convert(file, format, kpbs)
		if err != nil {
			h.Log.Warning(err.Error())
			user.AddMessage(fmt.Sprintf("Error converting file %s (%s). Try again with different parameters.", file.OriginalName, err.Error()))
		}
	}
	return nil
}
