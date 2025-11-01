package strategy

import (
	"rabbit-go/adb"
)

type RotationStrategy interface {
	Run() error
}

type RotationEnableStrategy struct{}

func (s *RotationEnableStrategy) Run() error {
	_, err := adb.Exec("adb shell settings put system accelerometer_rotation 1", false, nil)
	return err
}

type RotationDisableStrategy struct{}

func (s *RotationDisableStrategy) Run() error {
	_, err := adb.Exec("adb shell settings put system accelerometer_rotation 0", false, nil)
	return err
}

type RotationPortraitStrategy struct{}

func (s *RotationPortraitStrategy) Run() error {
	_, err := adb.Exec("adb shell settings put system user_rotation 0", false, nil)
	return err
}

type RotationLandscapeStrategy struct{}

func (s *RotationLandscapeStrategy) Run() error {
	_, err := adb.Exec("adb shell settings put system user_rotation 1", false, nil)
	return err
}

type RotationPortraitReverseStrategy struct{}

func (s *RotationPortraitReverseStrategy) Run() error {
	_, err := adb.Exec("adb shell settings put system user_rotation 2", false, nil)
	return err
}

type RotationLandscapeReverseStrategy struct{}

func (s *RotationLandscapeReverseStrategy) Run() error {
	_, err := adb.Exec("adb shell settings put system user_rotation 3", false, nil)
	return err
}
