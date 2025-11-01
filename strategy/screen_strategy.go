package strategy

import (
	"fmt"
	"rabbit-go/adb"
	"time"
)

type ScreenStrategy interface {
	Run() error
}

type ScreenshotStrategy struct{}

func (s *ScreenshotStrategy) Run() error {
	timestamp := time.Now().Format("2006_01_02_15_04_05")
	cmd := fmt.Sprintf("adb exec-out screencap -p > %s_screenshot.png", timestamp)
	_, err := adb.Exec(cmd, false, nil)
	return err
}

type Mp4RecordStrategy struct{}

func (s *Mp4RecordStrategy) Run() error {
	timestamp := time.Now().Format("2006_01_02_15_04_05")
	cmd := fmt.Sprintf("scrcpy --no-window -Nr %s_record.mp4", timestamp)
	_, err := adb.Exec(cmd, false, nil)
	return err
}
