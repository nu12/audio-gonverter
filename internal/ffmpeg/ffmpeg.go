package ffmpeg

import (
	"os"

	"github.com/nu12/audio-gonverter/internal/file"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Ffmpeg struct {
	InputPath  string
	OutputPath string
}

func (ff *Ffmpeg) Convert(f *file.File, format, kpbs string) error {

	convertedId := file.GenerateUUID()
	convertedName := f.Prefix() + "." + format

	if err := os.Mkdir(ff.OutputPath+"/"+convertedId, 0777); err != nil {

		return err
	}

	err := ffmpeg.Input(ff.InputPath+"/"+f.OriginalId+"/"+f.OriginalName).
		Output(ff.OutputPath+"/"+convertedId+"/"+convertedName, ffmpeg.KwArgs{"b:a": kpbs + "k"}).
		// OverWriteOutput().
		// ErrorToStdOut().
		Run()

	if err != nil {
		return err
	}

	f.ConvertedName = convertedName
	f.ConvertedId = convertedId
	f.IsConverted = true

	return nil
}
