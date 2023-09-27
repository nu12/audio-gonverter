package ffmpeg

import "github.com/nu12/audio-gonverter/internal/model"

type FfmpegMock struct{}

func (*FfmpegMock) Convert(file *model.File, format, kpbs string) error {
	return nil
}
