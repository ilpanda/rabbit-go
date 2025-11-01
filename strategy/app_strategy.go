package strategy

import (
	"fmt"
	"os"
	"path/filepath"
	"rabbit-go/adb"
	"rabbit-go/util"
	"strings"
)

type AppStrategy interface {
	CanHandle() bool
	Run(packageName string) error
	GetPackageName() string
}

// ClearAppDataStrategy clears app data
type ClearAppDataStrategy struct {
	PackageName string
}

func NewClearAppDataStrategy(packageName string) *ClearAppDataStrategy {
	return &ClearAppDataStrategy{PackageName: packageName}
}

func (s *ClearAppDataStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *ClearAppDataStrategy) Run(packageName string) error {
	cmd := fmt.Sprintf("adb shell pm clear %s", packageName)
	_, err := adb.Exec(cmd, false, nil)
	return err
}

func (s *ClearAppDataStrategy) GetPackageName() string {
	return s.PackageName
}

// KillStrategy force stops an app
type KillStrategy struct {
	PackageName string
}

func NewKillStrategy(packageName string) *KillStrategy {
	return &KillStrategy{PackageName: packageName}
}

func (s *KillStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *KillStrategy) Run(packageName string) error {
	cmd := fmt.Sprintf("adb shell am force-stop %s", packageName)
	_, err := adb.Exec(cmd, false, nil)
	return err
}

func (s *KillStrategy) GetPackageName() string {
	return s.PackageName
}

// GrantStrategy grants all permissions
type GrantStrategy struct {
	PackageName string
}

func NewGrantStrategy(packageName string) *GrantStrategy {
	return &GrantStrategy{PackageName: packageName}
}

func (s *GrantStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *GrantStrategy) Run(packageName string) error {
	cmd := fmt.Sprintf("adb shell dumpsys package %s", packageName)
	output, err := adb.Exec(cmd, false, nil)
	if err != nil {
		return err
	}

	permissions := getRequestedPermissions(util.MultiLine(output))
	for _, perm := range permissions {
		grantCmd := fmt.Sprintf("adb shell pm grant %s %s", packageName, perm)
		_, _ = adb.Exec(grantCmd, true, func(errorMsg string) bool {
			return strings.Contains(errorMsg, "Neither user 2000 nor current process has android.permission.GRANT_RUNTIME_PERMISSIONS")
		})
	}
	return nil
}

func (s *GrantStrategy) GetPackageName() string {
	return s.PackageName
}

func getRequestedPermissions(lines []string) []string {
	var permissions []string
	inPermissionSection := false

	for _, line := range lines {
		if !strings.Contains(line, ".permission.") {
			inPermissionSection = false
		}
		if strings.Contains(line, "requested permissions:") {
			inPermissionSection = true
			continue
		}
		if inPermissionSection {
			permissionName := strings.TrimSpace(strings.ReplaceAll(line, ":", ""))
			permissions = append(permissions, permissionName)
		}
	}
	return permissions
}

// RevokeStrategy revokes all permissions
type RevokeStrategy struct {
	PackageName string
}

func NewRevokeStrategy(packageName string) *RevokeStrategy {
	return &RevokeStrategy{PackageName: packageName}
}

func (s *RevokeStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *RevokeStrategy) Run(packageName string) error {
	cmd := fmt.Sprintf("adb shell dumpsys package %s", packageName)
	output, err := adb.Exec(cmd, false, nil)
	if err != nil {
		return err
	}

	lines := util.MultiLine(output)
	for _, line := range lines {
		if strings.Contains(line, "permission") && strings.Contains(line, "granted=true") {
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				permission := strings.TrimSpace(parts[0])
				revokeCmd := fmt.Sprintf("adb shell pm revoke %s %s", packageName, permission)
				_, _ = adb.Exec(revokeCmd, true, nil)
			}
		}
	}
	return nil
}

func (s *RevokeStrategy) GetPackageName() string {
	return s.PackageName
}

// StartActivityStrategy starts an app
type StartActivityStrategy struct {
	PackageName string
}

func NewStartActivityStrategy(packageName string) *StartActivityStrategy {
	return &StartActivityStrategy{PackageName: packageName}
}

func (s *StartActivityStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *StartActivityStrategy) Run(packageName string) error {
	cmd := fmt.Sprintf("adb shell monkey -p %s -c android.intent.category.LAUNCHER 1", packageName)
	_, err := adb.Exec(cmd, false, nil)
	return err
}

func (s *StartActivityStrategy) GetPackageName() string {
	return s.PackageName
}

// RestartAppStrategy restarts an app
type RestartAppStrategy struct {
	PackageName string
}

func NewRestartAppStrategy(packageName string) *RestartAppStrategy {
	return &RestartAppStrategy{PackageName: packageName}
}

func (s *RestartAppStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *RestartAppStrategy) Run(packageName string) error {
	killStrategy := NewKillStrategy(packageName)
	if err := killStrategy.Run(packageName); err != nil {
		return err
	}

	startStrategy := NewStartActivityStrategy(packageName)
	return startStrategy.Run(packageName)
}

func (s *RestartAppStrategy) GetPackageName() string {
	return s.PackageName
}

// StartAppDetailStrategy opens app details
type StartAppDetailStrategy struct {
	PackageName string
}

func NewStartAppDetailStrategy(packageName string) *StartAppDetailStrategy {
	return &StartAppDetailStrategy{PackageName: packageName}
}

func (s *StartAppDetailStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *StartAppDetailStrategy) Run(packageName string) error {
	cmd := fmt.Sprintf("adb shell am start -a android.settings.APPLICATION_DETAILS_SETTINGS package:%s", packageName)
	_, err := adb.Exec(cmd, false, nil)
	return err
}

func (s *StartAppDetailStrategy) GetPackageName() string {
	return s.PackageName
}

// ExportAppStrategy exports an app
type ExportAppStrategy struct {
	PackageName string
}

func NewExportAppStrategy(packageName string) *ExportAppStrategy {
	return &ExportAppStrategy{PackageName: packageName}
}

func (s *ExportAppStrategy) CanHandle() bool {
	return s.PackageName != ""
}

func (s *ExportAppStrategy) Run(packageName string) error {
	// Check if package exists
	cmd := fmt.Sprintf("adb shell pm list packages %s", packageName)
	output, err := adb.Exec(cmd, false, nil)
	if err != nil {
		return err
	}

	lines := util.MultiLine(output)
	packageExists := false
	for _, line := range lines {
		if strings.Contains(line, packageName) {
			packageExists = true
			break
		}
	}

	if !packageExists {
		util.Log(fmt.Sprintf("%s is not exists in phone", packageName))
		return nil
	}

	// Get APK path
	pathCmd := fmt.Sprintf("adb shell pm path %s", packageName)
	apkPath, err := adb.Exec(pathCmd, false, nil)
	if err != nil {
		return err
	}

	apkPath = strings.TrimPrefix(strings.TrimSpace(apkPath), "package:")
	if apkPath == "" {
		return fmt.Errorf("cannot find apk path")
	}

	// Pull APK
	destFile := fmt.Sprintf("./%s.apk", packageName)
	absPath, err := filepath.Abs(destFile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(absPath); err == nil {
		util.LogE(fmt.Sprintf("%s has exists", absPath))
		return nil
	}

	pullCmd := fmt.Sprintf("adb pull %s %s", apkPath, absPath)
	output, err = adb.Exec(pullCmd, false, nil)
	if err != nil {
		return err
	}

	util.Log(output)
	util.Log(fmt.Sprintf("apk has been saved in %s", absPath))
	return nil
}

func (s *ExportAppStrategy) GetPackageName() string {
	return s.PackageName
}
