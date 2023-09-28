package ffmpeg

import (
	"os"

	"github.com/nu12/audio-gonverter/internal/model"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Ffmpeg struct {
	InputPath  string
	OutputPath string
}

func (f *Ffmpeg) Convert(file *model.File, format, kpbs string) error {

	convertedId := model.GenerateUUID()
	convertedName := file.Prefix() + "." + format

	if err := os.Mkdir(f.OutputPath+convertedId, 0777); err != nil {

		return err
	}

	err := ffmpeg.Input(f.InputPath+file.OriginalId+"/"+file.OriginalName).
		Output(f.OutputPath+convertedId+"/"+convertedName, ffmpeg.KwArgs{"b:a": kpbs + "k"}).
		// OverWriteOutput().
		// ErrorToStdOut().
		Run()

	if err != nil {
		return err
	}

	file.ConvertedName = convertedName
	file.ConvertedId = convertedId
	file.IsConverted = true

	return nil
}
