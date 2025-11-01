package strategy

import (
	"fmt"
	"rabbit-go/adb"
	"rabbit-go/config"
	"rabbit-go/util"
)

type LogStrategy interface {
	CanHandle(packageName string, config config.LogConfig) bool
	Run(packageName string, config config.LogConfig) error
}

type LogCurrentActivityStrategy struct{}

func (s *LogCurrentActivityStrategy) CanHandle(packageName string, config config.LogConfig) bool {
	return config.LogCurrentActivity
}

func (s *LogCurrentActivityStrategy) Run(packageName string, config config.LogConfig) error {
	res, err := adb.GetCurrentPackageAndActivityName()
	if err != nil {
		return err
	}
	util.Log(res)
	return nil
}

type LogAllActivityStrategy struct{}

func (s *LogAllActivityStrategy) CanHandle(packageName string, config config.LogConfig) bool {
	return config.LogAllActivity
}

func (s *LogAllActivityStrategy) Run(packageName string, config config.LogConfig) error {
	res, err := adb.GetActivityListStringFromTopToBottom()
	if err != nil {
		return err
	}
	util.Log(res)
	return nil
}

type LogAllFragmentStrategy struct{}

func (s *LogAllFragmentStrategy) CanHandle(packageName string, config config.LogConfig) bool {
	return config.LogAllFragment
}

func (s *LogAllFragmentStrategy) Run(packageName string, config config.LogConfig) error {
	cmd := fmt.Sprintf(`adb shell dumpsys activity %s | grep -E '^\s*#\d' | grep -v -E 'ReportFragment|plan'`, packageName)
	res, err := adb.Exec(cmd, false, nil)
	if err != nil {
		return err
	}
	util.Log(res)
	return nil
}

type LogSpecificPackageActivityStrategy struct{}

func (s *LogSpecificPackageActivityStrategy) CanHandle(packageName string, config config.LogConfig) bool {
	return config.LogSpecificPackageActivity != ""
}

func (s *LogSpecificPackageActivityStrategy) Run(packageName string, config config.LogConfig) error {
	activityList, err := adb.GetActivityListStringFromTopToBottom()
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf(`echo '%s' | grep %s`, activityList, config.LogSpecificPackageActivity)
	res, err := adb.Exec(cmd, false, nil)
	if err != nil {
		return err
	}
	util.Log(res)
	return nil
}
