package adb

import (
	"rabbit-go/util"
	"strings"
)

func GetCurrentPackageAndActivityName() (string, error) {
	result, err := util.Exec(`adb shell dumpsys activity activities | grep mResumedActivity | awk '{print $4}'`, false, nil)

	if err != nil || strings.TrimSpace(result) == "" {
		result, err = util.Exec(`adb shell dumpsys activity activities | grep ResumedActivity | grep -v top | awk '{print $4}'`, false, nil)
		if err != nil {
			return "", err
		}
		return strings.TrimSuffix(result, "}\n"), nil
	}

	return strings.TrimSuffix(result, "}\n"), nil
}

func GetActivityListStringFromTopToBottom() (string, error) {
	return util.Exec(`adb shell dumpsys activity activities | grep -e 'Hist #' -e '* Hist'`, false, nil)
}

// Exec is a wrapper around util.Exec for convenience
func Exec(command string, ignoreError bool, exitWhen func(string) bool) (string, error) {
	return util.Exec(command, ignoreError, exitWhen)
}
