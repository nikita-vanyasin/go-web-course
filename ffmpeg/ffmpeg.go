package ffmpeg

import (
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetVideoDuration(videoPath string) (float64, error) {
	baseDir, err := filepath.Abs(path.Dir(videoPath))
	if err != nil {
		return 0.0, err
	}

	fileName := path.Base(videoPath)

	cmd := exec.Command(`docker`, `run`, `-t`, `--rm`, `-v`, baseDir+`:/files`, `sjourdan/ffprobe`, `-v`, `error`, `-show_entries`, `format=duration`, `-of`, `default=noprint_wrappers=1:nokey=1`, `/files/`+fileName)
	result, err := cmd.Output()
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(strings.Trim(string(result), "\n\r"), 64)
}

func ffmpegTimeFromSeconds(seconds int64) string {
	return time.Unix(seconds, 0).UTC().Format(`15:04:05.000000`)
}

func CreateVideoThumbnail(videoPath string, thumbnailPath string, thumbnailOffset int64) error {

	baseDir, err := filepath.Abs(path.Dir(videoPath))
	if err != nil {
		return err
	}

	fileName := path.Base(videoPath)
	thumbFileName := path.Base(thumbnailPath)

	return exec.Command(`docker`, `run`, `-t`, `--rm`, `-v`, baseDir+`:/files`, `sjourdan/ffmpeg`, `-i`, `/files/`+fileName, `-ss`, ffmpegTimeFromSeconds(thumbnailOffset), `-vframes`, `1`, `/files/`+thumbFileName).Run()
}
