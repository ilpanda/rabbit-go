package strategy

import (
	"fmt"
	"rabbit-go/adb"
	"rabbit-go/util"
	"strconv"
	"strings"
)

type DeviceInfoStrategy interface {
	Run() error
}

type DeviceInfoImpl struct{}

func (s *DeviceInfoImpl) Run() error {
	model, _ := adb.Exec("adb shell getprop ro.product.model", false, nil)
	version, _ := adb.Exec("adb shell getprop ro.build.version.release", false, nil)
	density, _ := adb.Exec("adb shell wm density", false, nil)
	display, _ := adb.Exec("adb shell dumpsys window displays", false, nil)
	androidID, _ := adb.Exec("adb shell settings get secure android_id", false, nil)
	sdkVersion, _ := adb.Exec("adb shell getprop ro.build.version.sdk", false, nil)
	ipAddress, _ := adb.Exec("adb shell ifconfig | grep Mask", true, nil)
	imei, _ := adb.Exec(`adb shell "service call iphonesubinfo 1 s16 com.android.shell | cut -c 52-66 | tr -d '.[:space:]'"`, false, nil)
	codeName, _ := adb.Exec("adb shell getprop ro.build.version.codename", false, nil)

	model = strings.TrimSpace(model)
	version = strings.TrimSpace(version)
	density = strings.TrimSpace(density)
	androidID = strings.TrimSpace(androidID)
	sdkVersion = strings.TrimSpace(sdkVersion)
	imei = strings.TrimSpace(imei)
	codeName = strings.ToUpper(strings.TrimSpace(codeName))

	if codeName == "REL" {
		codeName = ""
	}

	// Parse display
	displayLines := util.MultiLine(display)
	var displayRes string
	for _, line := range displayLines {
		if strings.Contains(line, "init=") {
			displayRes = strings.TrimSpace(line)
			if idx := strings.Index(displayRes, "rng"); idx != -1 {
				displayRes = displayRes[:idx]
			}
			break
		}
	}

	// Parse IP address
	ipAddressRes := ""
	permissionDeny := strings.Contains(ipAddress, "Permission denied")
	if !permissionDeny {
		ipAddressRes = fmt.Sprintf("ipAddress: %s", strings.TrimSpace(strings.ReplaceAll(ipAddress, "\n", "")))
	}

	// Parse density
	var densityRes string
	var densityScale float64
	var overrideRes string

	if !strings.Contains(density, "Override density") {
		idx := strings.Index(density, ":")
		if idx != -1 {
			densityRes = strings.TrimSpace(density[idx+1:])
			if d, err := strconv.ParseFloat(densityRes, 64); err == nil {
				densityScale = d / 160
			}
		}
	} else {
		lines := util.MultiLine(density)
		if len(lines) >= 2 {
			idx := strings.Index(lines[0], ":")
			if idx != -1 {
				densityRes = strings.TrimSpace(lines[0][idx+1:])
			}

			idx = strings.Index(lines[1], ":")
			if idx != -1 {
				overrideDensity := strings.TrimSpace(lines[1][idx+1:])
				if d, err := strconv.ParseFloat(overrideDensity, 64); err == nil {
					densityScale = d / 160
				}
				overrideRes = fmt.Sprintf("Override density: %sdpi", overrideDensity)
			}
		}
	}

	versionBuild := util.GetVersionBuild(sdkVersion)
	if versionBuild == "" {
		versionBuild = fmt.Sprintf("Android %s", version)
	}

	result := fmt.Sprintf(`model: %s
imei: %s
version: %s %s
display: %s
Physical density: %sdpi  %s
density scale: %.2f
android_id: %s
%s`,
		model,
		imei,
		versionBuild,
		codeName,
		displayRes,
		densityRes,
		overrideRes,
		densityScale,
		androidID,
		ipAddressRes,
	)

	util.Log(result)
	return nil
}

type CPUInfo struct{}

func (s *CPUInfo) Run() error {
	output, err := adb.Exec("adb shell cat /proc/cpuinfo", false, nil)
	if err != nil {
		return err
	}
	util.Log(output)
	return nil
}

type MemInfo struct{}

func (s *MemInfo) Run() error {
	output, err := adb.Exec("adb shell cat /proc/meminfo", false, nil)
	if err != nil {
		return err
	}
	util.Log(output)
	return nil
}

type BatteryInfo struct{}

func (s *BatteryInfo) Run() error {
	output, err := adb.Exec("adb shell dumpsys battery", false, nil)
	if err != nil {
		return err
	}
	util.Log(output)
	return nil
}
