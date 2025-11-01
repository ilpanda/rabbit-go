package cmd

import (
	"fmt"
	"os"
	"rabbit-go/adb"
	"rabbit-go/config"
	"rabbit-go/strategy"
	"strings"

	"github.com/spf13/cobra"
)

var (
	logConfig      config.LogConfig
	appConfig      config.AppConfig
	actionConfig   string
	infoConfig     string
	screenConfig   string
	rotationConfig string
)

var rootCmd = &cobra.Command{
	Use:   "rabbit-go",
	Short: "Android ADB command line tool",
	Long:  "A CLI tool for Android ADB operations",
}

func init() {
	// Log options
	rootCmd.Flags().BoolVarP(&logConfig.LogCurrentActivity, "current", "c", false, "print current activity name")
	rootCmd.Flags().BoolVarP(&logConfig.LogAllActivity, "all", "a", false, "print all activities name")
	rootCmd.Flags().BoolVarP(&logConfig.LogAllFragment, "fragment", "f", false, "print specific package fragments")
	rootCmd.Flags().StringVarP(&logConfig.LogSpecificPackageActivity, "print", "p", "", "print specific package activities")

	// App options
	rootCmd.Flags().StringVar(&appConfig.ClearAppPackageName, "clear", "", "clear app data")
	rootCmd.Flags().StringVar(&appConfig.KillAppPackageName, "kill", "", "force stop app")
	rootCmd.Flags().StringVar(&appConfig.GrantAppPermissionPackageName, "grant", "", "grant app all permissions")
	rootCmd.Flags().StringVar(&appConfig.RevokeAppPermissionPackageName, "revoke", "", "revoke app all permissions")
	rootCmd.Flags().StringVar(&appConfig.StartAppPackageName, "start", "", "start app")
	rootCmd.Flags().StringVar(&appConfig.RestartPackageName, "restart", "", "restart app")
	rootCmd.Flags().StringVar(&appConfig.StartAppDetailPackageName, "detail", "", "start app detail page")
	rootCmd.Flags().StringVar(&appConfig.ExportPackageName, "export", "", "export app to desktop")

	// Action config
	rootCmd.Flags().StringVar(&actionConfig, "action", "", "android adb start system activity (locale|developer|application|notification|bluetooth|input|display)")

	// Info config
	rootCmd.Flags().StringVarP(&infoConfig, "info", "i", "", "android adb get device info (device|cpu|memory|battery)")

	// Screen config
	rootCmd.Flags().StringVarP(&screenConfig, "screen", "s", "", "screenshot or record (png|mp4)")

	// Rotation config
	rootCmd.Flags().StringVarP(&rotationConfig, "rotate", "r", "", "screen rotation (enable|disable|0|1|2|3)")

	rootCmd.Run = runAdbCommand
}

func Execute() error {
	return rootCmd.Execute()
}

func runAdbCommand(cmd *cobra.Command, args []string) {
	res, err := adb.GetCurrentPackageAndActivityName()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current activity: %v\n", err)
		os.Exit(1)
	}

	parts := strings.Split(res, "/")
	if len(parts) == 0 {
		fmt.Fprintf(os.Stderr, "Invalid package/activity format\n")
		os.Exit(1)
	}

	packageName := strings.TrimSuffix(parts[0], "}")

	// Execute log commands
	executeLogCommands(packageName, logConfig)

	// Execute app commands
	executeAppCommands(appConfig)

	// Execute action config
	if actionConfig != "" {
		executeAction(actionConfig)
	}

	// Execute info config
	if infoConfig != "" {
		executeInfo(infoConfig)
	}

	// Execute screen config
	if screenConfig != "" {
		executeScreen(screenConfig)
	}

	// Execute rotation config
	if rotationConfig != "" {
		executeRotation(rotationConfig)
	}
}

func executeLogCommands(packageName string, config config.LogConfig) {
	strategies := []strategy.LogStrategy{
		&strategy.LogCurrentActivityStrategy{},
		&strategy.LogAllActivityStrategy{},
		&strategy.LogAllFragmentStrategy{},
		&strategy.LogSpecificPackageActivityStrategy{},
	}

	for _, s := range strategies {
		if s.CanHandle(packageName, config) {
			if err := s.Run(packageName, config); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
	}
}

func executeAppCommands(config config.AppConfig) {
	strategies := []strategy.AppStrategy{
		strategy.NewClearAppDataStrategy(config.ClearAppPackageName),
		strategy.NewKillStrategy(config.KillAppPackageName),
		strategy.NewGrantStrategy(config.GrantAppPermissionPackageName),
		strategy.NewRevokeStrategy(config.RevokeAppPermissionPackageName),
		strategy.NewStartActivityStrategy(config.StartAppPackageName),
		strategy.NewRestartAppStrategy(config.RestartPackageName),
		strategy.NewStartAppDetailStrategy(config.StartAppDetailPackageName),
		strategy.NewExportAppStrategy(config.ExportPackageName),
	}

	for _, s := range strategies {
		if s.CanHandle() {
			if err := s.Run(s.GetPackageName()); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
	}
}

func executeAction(action string) {
	actionMap := map[string]string{
		"locale":       "android.settings.LOCALE_SETTINGS",
		"developer":    "android.settings.APPLICATION_DEVELOPMENT_SETTINGS",
		"application":  "android.settings.APPLICATION_SETTINGS",
		"notification": "android.settings.ALL_APPS_NOTIFICATION_SETTINGS",
		"bluetooth":    "android.settings.BLUETOOTH_SETTINGS",
		"input":        "android.settings.INPUT_METHOD_SETTINGS",
		"display":      "android.settings.DISPLAY_SETTINGS",
	}

	if actionValue, ok := actionMap[action]; ok {
		cmd := fmt.Sprintf("adb shell am start -a %s", actionValue)
		if _, err := adb.Exec(cmd, false, nil); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing action: %v\n", err)
			os.Exit(1)
		}
	}
}

func executeInfo(info string) {
	var s strategy.DeviceInfoStrategy

	switch info {
	case "device":
		s = &strategy.DeviceInfoImpl{}
	case "cpu":
		s = &strategy.CPUInfo{}
	case "memory":
		s = &strategy.MemInfo{}
	case "battery":
		s = &strategy.BatteryInfo{}
	default:
		fmt.Fprintf(os.Stderr, "Unknown info type: %s\n", info)
		return
	}

	if err := s.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func executeScreen(screen string) {
	var s strategy.ScreenStrategy

	switch screen {
	case "png":
		s = &strategy.ScreenshotStrategy{}
	case "mp4":
		s = &strategy.Mp4RecordStrategy{}
	default:
		fmt.Fprintf(os.Stderr, "Unknown screen type: %s\n", screen)
		return
	}

	if err := s.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func executeRotation(rotation string) {
	var s strategy.RotationStrategy

	switch rotation {
	case "enable":
		s = &strategy.RotationEnableStrategy{}
	case "disable":
		s = &strategy.RotationDisableStrategy{}
	case "0":
		s = &strategy.RotationPortraitStrategy{}
	case "1":
		s = &strategy.RotationLandscapeStrategy{}
	case "2":
		s = &strategy.RotationPortraitReverseStrategy{}
	case "3":
		s = &strategy.RotationLandscapeReverseStrategy{}
	default:
		fmt.Fprintf(os.Stderr, "Unknown rotation type: %s\n", rotation)
		return
	}

	if err := s.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
