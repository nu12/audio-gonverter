package ffmpeg

import "github.com/nu12/audio-gonverter/internal/file"

type FfmpegMock struct{}

func (*FfmpegMock) Convert(file *file.File, format, kpbs string) error {
	return nil
}
